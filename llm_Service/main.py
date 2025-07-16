from fastapi import FastAPI

from task_Llm.upload_task_router import uploadTask

llm = FastAPI()
llm.include_router(uploadTask)



@llm.get("/")
def root():
    return {"Main page"}
    

