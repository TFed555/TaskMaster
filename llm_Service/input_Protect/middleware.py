import os
from jose import jwt
from jose.exceptions import JWTError, ExpiredSignatureError
from fastapi import HTTPException

class MiddleWare:
    secretKey = b"goida"
    def tokenVerification(self,token,userId):
        if not self.secretKey:
            HTTPException(status_code=401, detail="No secret key") 
        try:
            print("Python key:", self.secretKey)
            decoded_unverified = jwt.decode(token, key=self.secretKey, algorithms=["HS256"], options={"verify_signature": False})
            print("Unverified payload:", decoded_unverified)
        

        #     deCode = jwt.decode(token=token, key=self.secretKey, algorithms=["HS256"])
        #     if(str(deCode["sub"])!=str(userId)):#?
        #         raise HTTPException(status_code=401, detail="Invalid ") 
        #     return True
        except ExpiredSignatureError:
            raise HTTPException(status_code=401, detail="Expired token")
        except JWTError:
            raise HTTPException(status_code=401, detail="Invalid tok")
        except Exception as e:
            raise HTTPException(status_code=401, detail=f"Token verification failed: {str(e)}")
    
        # decoded_unverified = jwt.decode(token=str(token), algorithms="HS256", key=str(self.secretKey), options={"verify_signature": False})
        # print(decoded_unverified)
    #     decoded = jwt.decode(
    #         token,
    #         self.secretKey,
    #         algorithms=["HS256"],
    #         options={"verify_exp": True}
    #     )
        
    #     if decoded.get("sub") != userId:
    #         raise HTTPException(status_code=401, detail="Invalid token") 
        
    #     return True
    
    # except jwt.ExpiredSignatureError:
    #     raise HTTPException(status_code=401, detail="Expired token")
    # except jwt.InvalidTokenError:
    #     raise HTTPException(status_code=401, detail="Invalid token")
    # except Exception as e:
    #     raise HTTPException(status_code=401, detail=f"Token verification failed: {str(e)}")