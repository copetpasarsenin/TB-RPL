from fastapi import FastAPI, Request, HTTPException, Depends, Header, Response
from fastapi.responses import JSONResponse
from fastapi.staticfiles import StaticFiles
from sqlalchemy import create_engine, Column, Integer, String, Float, DateTime
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, Session
import jwt
import httpx
import datetime
import os
from dotenv import load_dotenv

load_dotenv()

app = FastAPI(
    title="API Gateway - Final (Kelompok 3)",
    description="API Gateway untuk ekosistem UMKM (Ported to Python)",
    version="1.0.0"
)

# --- DATABASE SETUP ---
DB_USER = os.getenv("DB_USER", "postgres")
DB_PASSWORD = os.getenv("DB_PASSWORD", "admin")
DB_HOST = os.getenv("DB_HOST", "localhost")
DB_PORT = os.getenv("DB_PORT", "5432")
DB_NAME = os.getenv("DB_NAME", "postgres")

SQLALCHEMY_DATABASE_URL = f"postgresql://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/{DB_NAME}"
engine = create_engine(SQLALCHEMY_DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()

# --- MODELS ---
class RequestLog(Base):
    __tablename__ = "request_logs_py"
    id = Column(Integer, primary_key=True, index=True)
    method = Column(String)
    path = Column(String)
    source_ip = Column(String)
    status_code = Column(Integer)
    created_at = Column(DateTime, default=datetime.datetime.utcnow)

class FeeTransaction(Base):
    __tablename__ = "fee_transactions_py"
    id = Column(Integer, primary_key=True, index=True)
    original_amount = Column(Float)
    fee_percent = Column(Float)
    fee_amount = Column(Float)
    source_service = Column(String)
    destination_service = Column(String)
    created_at = Column(DateTime, default=datetime.datetime.utcnow)

Base.metadata.create_all(bind=engine)

def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

# --- GLOBAL LOGGING MIDDLEWARE (Aturan #42) ---
# Mencatat SEMUA aktivitas secara otomatis ke Database
@app.middleware("http")
async def log_requests(request: Request, call_next):
    db = SessionLocal()
    status_code = 500
    try:
        response = await call_next(request)
        status_code = response.status_code
        return response
    except Exception as e:
        raise e
    finally:
        # Hindari nge-log halaman UI dan Dokumntasi biar ga penuh
        ignored_paths = ["/ui", "/docs", "/openapi.json"]
        if not any(request.url.path.startswith(p) for p in ignored_paths):
            log = RequestLog(
                method=request.method,
                path=request.url.path,
                source_ip=request.client.host if request.client else "127.0.0.1",
                status_code=status_code
            )
            db.add(log)
            db.commit()
        db.close()

# --- SECURITY ---
SECRET_KEY = "supersecretkey"

def verify_token(authorization: str = Header(None)):
    if not authorization or not authorization.startswith("Bearer "):
        raise HTTPException(status_code=401, detail={"status": "error", "message": "Authorization header diperlukan"})
    token = authorization.split(" ")[1]
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=["HS256"])
        return payload
    except jwt.ExpiredSignatureError:
        raise HTTPException(status_code=401, detail={"status": "error", "message": "Token kadaluarsa"})
    except jwt.InvalidTokenError:
        raise HTTPException(status_code=401, detail={"status": "error", "message": "Token tidak valid"})

# --- ROUTES SESUAI SPREADSHEET ---

@app.get("/health", summary="Health Check Gateway")
def health_check():
    return {"status": "success", "message": "API Gateway is running (Python Version)", "service": "api-gateway"}

@app.post("/auth/token", summary="Generate JWT Token (Dummy)")
def generate_token(payload: dict):
    token = jwt.encode({**payload, "exp": datetime.datetime.utcnow() + datetime.timedelta(hours=24)}, SECRET_KEY, algorithm="HS256")
    return {"status": "success", "data": {"token": token}}

@app.post("/integrator/validasi_request", summary="Validasi Request (Token)")
def validasi_request(user=Depends(verify_token)):
    return {"status": "success", "data": "Token valid", "user": user}

@app.get("/integrator/logging", summary="Logging (Mencatat Request)")
def get_logs(db: Session = Depends(get_db), user=Depends(verify_token)):
    logs = db.query(RequestLog).order_by(RequestLog.created_at.desc()).limit(100).all()
    return {"status": "success", "data": logs}

@app.post("/integrator/biaya_layanan_integrasi", summary="Biaya Layanan Integrasi (Potong Fee 0.5%)")
def biaya_layanan(payload: dict, db: Session = Depends(get_db), user=Depends(verify_token)):
    amount = payload.get("amount", 0)
    fee_percent = 0.5
    fee_amount = amount * (fee_percent / 100)
    
    fee = FeeTransaction(
        original_amount=amount,
        fee_percent=fee_percent,
        fee_amount=fee_amount,
        source_service=payload.get("source_service", "unknown"),
        destination_service=payload.get("destination_service", "unknown")
    )
    db.add(fee)
    db.commit()
    db.refresh(fee)
    
    return {
        "status": "success",
        "data": {
            "transaction_id": fee.id,
            "original_amount": amount,
            "fee_percent": fee_percent,
            "fee_amount": fee_amount,
            "total_after_fee": amount + fee_amount
        }
    }

SERVICE_REGISTRY = {
    "smartbank": "http://localhost:8081",
    "marketplace": "http://localhost:8082",
    "pos": "http://localhost:8083",
    "supplierhub": "http://localhost:8084",
    "logistikita": "http://localhost:8085",
    "umkm_insight": "http://localhost:8086"
}

@app.api_route("/integrator/routing_api/{service_name}/{path:path}", methods=["GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"], summary="Routing API (Forward Request)")
async def proxy_to_service(request: Request, service_name: str, path: str, db: Session = Depends(get_db)):
    auth_header = request.headers.get("authorization")
    if request.method != "OPTIONS":
        if not auth_header or not auth_header.startswith("Bearer "):
            return JSONResponse(status_code=401, content={"status": "error", "message": "Authorization header diperlukan"})
        token = auth_header.split(" ")[1]
        try:
            jwt.decode(token, SECRET_KEY, algorithms=["HS256"])
        except Exception:
            return JSONResponse(status_code=401, content={"status": "error", "message": "Token tidak valid"})

    target_base = SERVICE_REGISTRY.get(service_name)
    if not target_base:
        return JSONResponse(status_code=404, content={"status": "error", "message": "Service tidak ditemukan dalam Registry Gateway"})

    target_url = f"{target_base}/{path}"
    body = await request.body()
    headers = dict(request.headers)
    headers.pop("host", None) 
    
    status_code = 500
    try:
        async with httpx.AsyncClient() as client:
            resp = await client.request(
                method=request.method,
                url=target_url,
                headers=headers,
                content=body,
                params=request.query_params
            )
            status_code = resp.status_code
            response_content = resp.content
            response_headers = dict(resp.headers)
            response_headers.pop('content-encoding', None)
            response_headers.pop('content-length', None)
            response_headers.pop('transfer-encoding', None)
    except httpx.RequestError:
        status_code = 502
        response_content = b'{"status": "error", "message": "Gagal menghubungi service tujuan"}'
        response_headers = {"content-type": "application/json"}

    return Response(content=response_content, status_code=status_code, headers=response_headers)

app.mount("/ui", StaticFiles(directory="public", html=True), name="ui")
