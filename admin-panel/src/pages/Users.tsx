import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../api/client';
import { routes } from '../routes';
import { formatDate } from '../utils/format';

export function UsersPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const [searchInput, setSearchInput] = useState('');
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['users', page, search],
    queryFn: () => api.listUsers(page, 20, search),
  });

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this user?')) return;
    await api.deleteUser(id);
    queryClient.invalidateQueries({ queryKey: ['users'] });
  };

  return (
    <div>
      <h1 className="page-title">Users</h1>
      <div className="toolbar">
        <input
          placeholder="Search by email or telegram ID..."
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          style={{ maxWidth: 320 }}
        />
        <button
          type="button"
          className="btn btn-primary"
          onClick={() => {
            setPage(1);
            setSearch(searchInput);
          }}
        >
          Search
        </button>
      </div>

      <div className="card table-wrap">
        {isLoading ? (
          <div className="loading">Loading...</div>
        ) : (
          <table>
            <thead>
              <tr>
                <th>Telegram ID</th>
                <th>Email</th>
                <th>User ID</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {data?.users.map((user) => (
                <tr key={user.userID}>
                  <td>{user.telegramID || '—'}</td>
                  <td>{user.email || '—'}</td>
                  <td>
                    <code>{user.userID.slice(0, 8)}...</code>
                  </td>
                  <td>{formatDate(user.createdAt)}</td>
                  <td className="actions">
                    <Link to={routes.user(user.userID)} className="btn btn-secondary btn-sm">
                      View
                    </Link>
                    <button
                      type="button"
                      className="btn btn-danger btn-sm"
                      onClick={() => handleDelete(user.userID)}
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
        <button
          type="button"
          className="btn btn-secondary btn-sm"
          disabled={page <= 1}
          onClick={() => setPage((p) => p - 1)}
        >
          Prev
        </button>
        <span>
          Page {page} of {Math.max(1, Math.ceil((data?.total ?? 0) / 20))}
        </span>
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
