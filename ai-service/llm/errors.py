from config.prompts import LLM_FALLBACK_ERROR_MESSAGE as G4F_FALLBACK_ERROR_MESSAGE

class LLMUserFacingError(Exception):
    """Ошибка с текстом, безопасным для показа пользователю."""

    def __init__(self, message: str):
        self.message = message
        super().__init__(message)
