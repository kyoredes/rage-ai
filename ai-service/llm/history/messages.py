from typing import Literal

from langchain_core.messages import AIMessage, BaseMessage, HumanMessage, SystemMessage
from pydantic import BaseModel

from config.prompts import G4F_INSTRUCTION_ACK


class ChatMessage(BaseModel):
    role: Literal["user", "assistant"]
    content: str


def to_openai_messages(
    system_prompt: str,
    messages: list[ChatMessage],
) -> list[dict[str, str]]:
    return [
        {"role": "system", "content": system_prompt},
        *[message.model_dump() for message in messages],
    ]


def to_g4f_messages(
    system_prompt: str,
    messages: list[ChatMessage],
) -> list[dict[str, str]]:
    """g4f free providers often ignore role=system; use a user/assistant pair instead."""
    result: list[dict[str, str]] = [
        {"role": "user", "content": system_prompt},
        {"role": "assistant", "content": G4F_INSTRUCTION_ACK},
    ]

    dumped = [message.model_dump() for message in messages]
    has_prior_assistant = any(item["role"] == "assistant" for item in dumped[:-1])

    if dumped and dumped[-1]["role"] == "user" and has_prior_assistant:
        dumped[-1] = {
            "role": "user",
            "content": f"[Инструкция: {system_prompt}]\n\n{dumped[-1]['content']}",
        }

    result.extend(dumped)
    return result


def to_langchain_messages(
    system_prompt: str,
    messages: list[ChatMessage],
) -> list[BaseMessage]:
    result: list[BaseMessage] = [SystemMessage(content=system_prompt)]
    for message in messages:
        if message.role == "user":
            result.append(HumanMessage(content=message.content))
        else:
            result.append(AIMessage(content=message.content))
    return result
