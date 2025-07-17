import os
from jose import jwt
from fastapi import HTTPException

class MiddleWare:
    secretKey = "goida"
    def tokenVerification(self,token,userId):
        if self.secretKey=="":
            HTTPException(status_code=401, detail="No secret key") 
        try:
            deCode = jwt.decode(token=token, key=self.secretKey, algorithms="HS256")
            if(deCode["sub"]!=userId):#?
                raise HTTPException(status_code=401, detail="Invalid token") 
            return True
        except Exception:
            print(type(token))
            print(self.secretKey)
           # print(deCode)
            raise HTTPException(status_code=401, detail="Expired token")