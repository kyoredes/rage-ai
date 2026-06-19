from aiogram import Router, F, types
from aiogram.types import Message, ReplyKeyboardRemove
from aiogram.fsm.context import FSMContext
from core.keyboards import get_start_keyboard

main_router = Router()


@main_router.message(F.text == "👈 Назад")
async def cancel(message: Message, state: FSMContext):
    await state.clear()
    await message.answer("Отправляемся назад", reply_markup=ReplyKeyboardRemove())
    await message.answer(
        "🏠 Вы вернулись в главное меню", reply_markup=get_start_keyboard()
    )


@main_router.callback_query(F.data == "main_menu")
async def back_to_main_menu(callback: types.CallbackQuery, state: FSMContext):
    await state.clear()
    await callback.message.answer(
        "🏠 Вы вернулись в главное меню", reply_markup=get_start_keyboard()
    )
    await callback.answer()
