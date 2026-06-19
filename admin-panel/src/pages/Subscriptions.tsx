import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { api } from '../api/client';
import { formatDate, fromDatetimeLocal, toDatetimeLocal } from '../utils/format';

export function SubscriptionsPage() {
  const [page, setPage] = useState(1);
  const [status, setStatus] = useState('');
  const [editingId, setEditingId] = useState<string | null>(null);
  const [startsAt, setStartsAt] = useState('');
  const [expiresAt, setExpiresAt] = useState('');
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['subscriptions', page, status],
    queryFn: () => api.listSubscriptions(page, 20, status),
  });

  const updateMutation = useMutation({
    mutationFn: () =>
      api.updateSubscription(editingId!, fromDatetimeLocal(startsAt), fromDatetimeLocal(expiresAt)),
    onSuccess: () => {
      setEditingId(null);
      queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
    },
  });

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this subscription?')) return;
    await api.deleteSubscription(id);
    queryClient.invalidateQueries({ queryKey: ['subscriptions'] });
  };

  const startEdit = (sub: { subscriptionID: string; startsAt: number; expiresAt: number }) => {
    setEditingId(sub.subscriptionID);
    setStartsAt(toDatetimeLocal(sub.startsAt));
    setExpiresAt(toDatetimeLocal(sub.expiresAt));
  };

  const isActive = (expiresAt: number) => expiresAt * 1000 > Date.now();

  return (
    <div>
      <h1 className="page-title">Subscriptions</h1>
      <div className="toolbar">
        <select value={status} onChange={(e) => { setStatus(e.target.value); setPage(1); }} style={{ maxWidth: 200 }}>
          <option value="">All</option>
          <option value="active">Active</option>
          <option value="expired">Expired</option>
        </select>
      </div>

      {editingId && (
        <div className="card" style={{ marginBottom: '1rem' }}>
          <h3>Edit Subscription</h3>
          <div className="toolbar">
            <div className="form-group" style={{ flex: 1 }}>
              <label>Starts At</label>
              <input type="datetime-local" value={startsAt} onChange={(e) => setStartsAt(e.target.value)} />
            </div>
            <div className="form-group" style={{ flex: 1 }}>
              <label>Expires At</label>
              <input type="datetime-local" value={expiresAt} onChange={(e) => setExpiresAt(e.target.value)} />
            </div>
          </div>
          <div className="actions">
            <button type="button" className="btn btn-primary btn-sm" onClick={() => updateMutation.mutate()}>
              Save
            </button>
            <button type="button" className="btn btn-secondary btn-sm" onClick={() => setEditingId(null)}>
              Cancel
            </button>
          </div>
        </div>
      )}

      <div className="card table-wrap">
        {isLoading ? (
          <div className="loading">Loading...</div>
        ) : (
          <table>
            <thead>
              <tr>
                <th>User ID</th>
                <th>Starts</th>
                <th>Expires</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {data?.subscriptions.map((sub) => (
                <tr key={sub.subscriptionID}>
                  <td>
                    <code>{sub.userID.slice(0, 8)}...</code>
                  </td>
                  <td>{formatDate(sub.startsAt)}</td>
                  <td>{formatDate(sub.expiresAt)}</td>
                  <td>
                    <span className={`badge ${isActive(sub.expiresAt) ? 'badge-active' : 'badge-expired'}`}>
                      {isActive(sub.expiresAt) ? 'Active' : 'Expired'}
                    </span>
                  </td>
                  <td className="actions">
                    <button type="button" className="btn btn-secondary btn-sm" onClick={() => startEdit(sub)}>
                      Edit
                    </button>
                    <button
                      type="button"
                      className="btn btn-danger btn-sm"
                      onClick={() => handleDelete(sub.subscriptionID)}
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      <div className="pagination">
        <button type="button" className="btn btn-secondary btn-sm" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>
          Prev
        </button>
        <span>Page {page}</span>
        <button
          type="button"
          className="btn btn-secondary btn-sm"
          disabled={!data || page * 20 >= data.total}
          onClick={() => setPage((p) => p + 1)}
        >
          Next
        </button>
      </div>
    </div>
  );
}
