import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { api } from '../api/client';

export function ChatHistoryPage() {
  const [page, setPage] = useState(1);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const queryClient = useQueryClient();

  const { data: sessions, isLoading } = useQuery({
    queryKey: ['chat-sessions', page],
    queryFn: () => api.listChatSessions(page, 20),
  });

  const { data: history, isLoading: historyLoading } = useQuery({
    queryKey: ['chat-history', selectedId],
    queryFn: async () => (await api.getChatHistory(selectedId!)).history,
    enabled: !!selectedId,
  });

  const handleClear = async (telegramId: string) => {
    if (!confirm('Clear chat history for this session?')) return;
    await api.clearChatHistory(telegramId);
    queryClient.invalidateQueries({ queryKey: ['chat-sessions'] });
    queryClient.invalidateQueries({ queryKey: ['chat-history', telegramId] });
    if (selectedId === telegramId) setSelectedId(null);
  };

  return (
    <div>
      <h1 className="page-title">Chat History</h1>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
        <div className="card table-wrap">
          <h3>Sessions</h3>
          {isLoading ? (
            <div className="loading">Loading...</div>
          ) : (
            <table>
              <thead>
                <tr>
                  <th>Telegram ID</th>
                  <th>Messages</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {sessions?.sessions.map((s) => (
                  <tr key={s.telegramID}>
                    <td>{s.telegramID}</td>
                    <td>{s.messageCount}</td>
                    <td className="actions">
                      <button
                        type="button"
                        className="btn btn-secondary btn-sm"
                        onClick={() => setSelectedId(s.telegramID)}
                      >
                        View
                      </button>
                      <button
                        type="button"
                        className="btn btn-danger btn-sm"
                        onClick={() => handleClear(s.telegramID)}
                      >
                        Clear
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
          <div className="pagination">
            <button type="button" className="btn btn-secondary btn-sm" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>
              Prev
            </button>
            <span>Page {page}</span>
            <button
              type="button"
              className="btn btn-secondary btn-sm"
              disabled={!sessions || page * 20 >= sessions.total}
              onClick={() => setPage((p) => p + 1)}
            >
              Next
            </button>
          </div>
        </div>

        <div className="card">
          <h3>Messages {selectedId && `— ${selectedId}`}</h3>
          {!selectedId && <div className="loading">Select a session</div>}
          {selectedId && historyLoading && <div className="loading">Loading...</div>}
          {history?.messages.map((msg, i) => (
            <div key={i} className={`message message-${msg.role}`}>
              <div className="message-role">{msg.role}</div>
              <div>{msg.content}</div>
            </div>
          ))}
          {selectedId && history && history.messages.length === 0 && (
            <div className="loading">No messages</div>
          )}
        </div>
      </div>
    </div>
  );
}
