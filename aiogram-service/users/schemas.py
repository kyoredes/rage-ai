from pydantic import BaseModel


class UserModel(BaseModel):
    status: str
