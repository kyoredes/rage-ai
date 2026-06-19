import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { api } from '../api/client';
import { routes } from '../routes';
import { formatDate } from '../utils/format';

export function UserDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [email, setEmail] = useState('');
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['user', id],
    queryFn: async () => (await api.getUser(id!)).user,
    enabled: !!id,
  });

  useEffect(() => {
    if (data) setEmail(data.email);
  }, [data]);

  const updateMutation = useMutation({
    mutationFn: () => api.updateUser(id!, email),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user', id] });
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });

  if (isLoading) return <div className="loading">Loading...</div>;
  if (!data) return <div className="error">User not found</div>;

  return (
    <div>
      <Link to={routes.users} style={{ color: 'var(--muted)', fontSize: '0.875rem' }}>
        ← Back to users
      </Link>
      <h1 className="page-title">User Detail</h1>

      <div className="card detail-grid">
        <div>
          <strong>User ID:</strong> {data.userID}
        </div>
        <div>
          <strong>Telegram ID:</strong> {data.telegramID || '—'}
        </div>
        <div>
          <strong>Created:</strong> {formatDate(data.createdAt)}
        </div>
        <div>
          <strong>Updated:</strong> {formatDate(data.updatedAt)}
        </div>

        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input id="email" value={email} onChange={(e) => setEmail(e.target.value)} />
        </div>
        <button
          type="button"
          className="btn btn-primary"
          onClick={() => updateMutation.mutate()}
          disabled={updateMutation.isPending}
        >
          {updateMutation.isPending ? 'Saving...' : 'Update Email'}
        </button>
      </div>
    </div>
  );
}
