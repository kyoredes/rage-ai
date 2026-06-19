import asyncio
from contextlib import asynccontextmanager

from aiogram.enums import ChatAction
from aiogram.types import Message

TYPING_REFRESH_SECONDS = 4


@asynccontextmanager
async def show_typing(message: Message):
    stop = asyncio.Event()

    async def refresh_typing() -> None:
        while not stop.is_set():
            try:
                await message.bot.send_chat_action(
                    chat_id=message.chat.id,
                    action=ChatAction.TYPING,
                )
            except Exception:
                return

            try:
                await asyncio.wait_for(stop.wait(), timeout=TYPING_REFRESH_SECONDS)
            except TimeoutError:
                continue

    task = asyncio.create_task(refresh_typing())
    try:
        yield
    finally:
        stop.set()
        task.cancel()
        try:
            await task
        except asyncio.CancelledError:
            pass
