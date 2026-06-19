from config.prompts import ASSISTANT_SYSTEM_PROMPT, G4F_SYSTEM_PROMPT
from config.settings import settings
from redis.asyncio import Redis

PROMPT_KEY = "llm:system_prompt"


def default_system_prompt(*, for_g4f: bool = False) -> str:
    return G4F_SYSTEM_PROMPT if for_g4f else ASSISTANT_SYSTEM_PROMPT


class SystemPromptStore:
    def __init__(self, redis_client: Redis | None = None):
        self._redis = redis_client or Redis.from_url(
            settings.REDIS_URL,
            decode_responses=True,
            socket_timeout=5,
            socket_connect_timeout=5,
        )

    async def get_custom(self) -> str | None:
        value = await self._redis.get(PROMPT_KEY)
        if value is None:
            return None
        stripped = value.strip()
        return stripped if stripped else None

    async def get_effective(self, *, for_g4f: bool = False) -> str:
        custom = await self.get_custom()
        if custom is not None:
            return custom
        return default_system_prompt(for_g4f=for_g4f)

    async def set_custom(self, prompt: str) -> None:
        stripped = prompt.strip()
        if not stripped:
            await self._redis.delete(PROMPT_KEY)
            return
        await self._redis.set(PROMPT_KEY, stripped)

    async def get_admin_view(self) -> dict:
        custom = await self.get_custom()
        return {
            "prompt": custom if custom is not None else ASSISTANT_SYSTEM_PROMPT,
            "default_prompt": ASSISTANT_SYSTEM_PROMPT,
            "is_custom": custom is not None,
        }
