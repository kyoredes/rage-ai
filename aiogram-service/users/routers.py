from aiogram import Router
from aiogram.filters import CommandStart
from aiogram.types import Message
from core.keyboards import get_start_keyboard
from users.manager import UserManager

users_router = Router(name="users")

@users_router.message(CommandStart())
async def start(message: Message):
    user_manager = UserManager()
    client = await user_manager.start(str(message.from_user.id))
    if client is None:
        await message.answer("Не удалось создать клиента")
        return

    await message.answer("Привет! Выбери действие:", reply_markup=get_start_keyboard())
