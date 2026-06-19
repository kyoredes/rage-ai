import { useQuery } from '@tanstack/react-query';
import { api, type ServiceStatus } from '../api/client';
import { formatDate } from '../utils/format';

function statusClass(status: ServiceStatus['status']): string {
  if (status === 'up') return 'status-up';
  if (status === 'degraded') return 'status-degraded';
  return 'status-down';
}

function statusLabel(status: ServiceStatus['status']): string {
  if (status === 'up') return 'Up';
  if (status === 'degraded') return 'Degraded';
  return 'Down';
}

export function ServicesStatusPanel() {
  const { data, isLoading, error, dataUpdatedAt } = useQuery({
    queryKey: ['services-status'],
    queryFn: async () => (await api.getServicesStatus()).servicesStatus,
    refetchInterval: 30_000,
  });

  if (isLoading) return <div className="loading">Checking services...</div>;
  if (error) return <div className="error">Failed to load service statuses</div>;

  return (
    <div className="card" style={{ marginBottom: '1.5rem' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
        <h3 style={{ margin: 0 }}>Services</h3>
        <span style={{ color: 'var(--muted)', fontSize: '0.8rem' }}>
          Updated {formatDate(Math.floor(dataUpdatedAt / 1000))}
        </span>
      </div>
      <div className="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Service</th>
              <th>Status</th>
              <th>Latency</th>
            </tr>
          </thead>
          <tbody>
            {data?.services.map((service) => (
              <tr key={service.id}>
                <td>{service.name}</td>
                <td>
                  <span className={`status-badge ${statusClass(service.status)}`}>
                    {statusLabel(service.status)}
                  </span>
                </td>
                <td>
                  {service.latencyMs === 0 && service.id === 'gateway'
                    ? '<1 ms'
                    : `${service.latencyMs} ms`}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
