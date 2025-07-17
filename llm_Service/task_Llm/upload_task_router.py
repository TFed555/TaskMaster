from task_Llm.task_master_llm import TaskMasterLLM
from fastapi import APIRouter
import os
from jose import jwt
from input_Protect.validation_input import ReqTaskLLM
from input_Protect.middleware import MiddleWare
import requests

taskMasterLLM = TaskMasterLLM()
uploadTask = APIRouter(prefix = "/api/task_llm",tags=["Отправка цели и описание цели в LLM, для генерации шагов достижения"])



@uploadTask.post("/postTask")
def  sendToChequeInfo(reqTaskLLM: ReqTaskLLM):
    urlHost = "https://localhost:8082"
    plan = taskMasterLLM.getPlan(task=reqTaskLLM.task,description=reqTaskLLM.description,todoId = reqTaskLLM.todoID)
    url = "http://notes_service:8082/api/todos/plan"
    requests.post(url=url, data=plan)
    return "Send successfully"
 
