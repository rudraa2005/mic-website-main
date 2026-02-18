from fastapi import FastAPI, HTTPException, Request
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import Optional
import os
import sys

from Insights import analyze_startup_idea
from Chatbot import chat_reply

app = FastAPI()

# Allow Go backend to proxy requests
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

class AnalyzeRequest(BaseModel):
    submission_id: str
    file_path: str  # ABSOLUTE PATH from Go

class ChatRequest(BaseModel):
    message: str
    session_id: Optional[str] = None

@app.get("/")
def health():
    return {"status": "AI service running"}

@app.post("/chat")
def chat(req: ChatRequest, request: Request):
    """Chatbot endpoint — proxied from Go via /api/chat"""
    session_id = req.session_id or request.headers.get("X-Session-ID")
    result = chat_reply(req.message, session_id)
    return result

@app.post("/analyze", status_code=202)
def analyze(req: AnalyzeRequest):
    print("[PYTHON] /analyze HIT", file=sys.stderr)
    print("[PYTHON] submission_id:", req.submission_id, file=sys.stderr)
    print("[PYTHON] file_path:", req.file_path, file=sys.stderr)

    full_path = os.path.normpath(req.file_path)

    print("[PYTHON] resolved full_path:", full_path, file=sys.stderr)
    print("[PYTHON] file exists:", os.path.exists(full_path), file=sys.stderr)

    if not os.path.exists(full_path):
        raise HTTPException(
            status_code=404,
            detail=f"File not found at {full_path}",
        )

    result = analyze_startup_idea(full_path)

    if not result:
        raise HTTPException(status_code=500, detail="Analysis failed")

    if "error" in result:
        raise HTTPException(status_code=502, detail=result)

    return {
        "submission_id": req.submission_id,
        "viability": result["viability"],
        "fit": result["fit"],
        "saturation": result["saturation"],
        "recommendations": result["recommendations"],
    }