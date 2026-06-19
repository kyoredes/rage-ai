import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useEffect, useState } from 'react';
import { api } from '../api/client';

export function LLMConfigPage() {
  const [prompt, setPrompt] = useState('');
  const queryClient = useQueryClient();

  const { data: config, isLoading: configLoading, error: configError } = useQuery({
    queryKey: ['llm-config'],
    queryFn: async () => (await api.getLLMConfig()).config,
  });

  const { data: systemPrompt, isLoading: promptLoading, error: promptError } = useQuery({
    queryKey: ['system-prompt'],
    queryFn: async () => (await api.getSystemPrompt()).systemPrompt,
  });

  useEffect(() => {
    if (systemPrompt) setPrompt(systemPrompt.prompt);
  }, [systemPrompt]);

  const saveMutation = useMutation({
    mutationFn: () => api.updateSystemPrompt(prompt),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['system-prompt'] });
    },
  });

  const resetMutation = useMutation({
    mutationFn: () => api.updateSystemPrompt(''),
    onSuccess: (result) => {
      setPrompt(result.systemPrompt.prompt);
      queryClient.invalidateQueries({ queryKey: ['system-prompt'] });
    },
  });

  if (configLoading || promptLoading) return <div className="loading">Loading...</div>;
  if (configError || promptError) return <div className="error">Failed to load config</div>;

  return (
    <div>
      <h1 className="page-title">LLM Configuration</h1>

      <div className="card detail-grid" style={{ marginBottom: '1.5rem' }}>
        <div>
          <strong>Provider:</strong> {config?.provider}
        </div>
        <div>
          <strong>Debug Mode:</strong>{' '}
          <span className={`badge ${config?.debug ? 'badge-active' : 'badge-expired'}`}>
            {config?.debug ? 'ON (G4F only)' : 'OFF (LiteLLM + G4F)'}
          </span>
        </div>
        {config?.usesLiteLLM && (
          <>
            <div>
              <strong>LiteLLM Model:</strong> {config.model}
            </div>
            <div>
              <strong>Temperature:</strong> {config.temperature}
            </div>
            <div>
              <strong>Max Tokens:</strong> {config.maxTokens}
            </div>
          </>
        )}
        {config?.g4fModels && config.g4fModels.length > 0 && (
          <div style={{ gridColumn: '1 / -1' }}>
            <strong>G4F Models{config.usesLiteLLM ? ' (fallback)' : ''}:</strong>
            <ol style={{ margin: '0.5rem 0 0', paddingLeft: '1.25rem' }}>
              {config.g4fModels.map((model) => (
                <li key={model}>{model}</li>
              ))}
            </ol>
          </div>
        )}
      </div>

      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
          <h3 style={{ margin: 0 }}>System Prompt</h3>
          {systemPrompt?.isCustom ? (
            <span className="badge badge-active">Custom</span>
          ) : (
            <span className="badge badge-expired">Default</span>
          )}
        </div>

        <div className="form-group">
          <label htmlFor="system-prompt">Prompt for the neural network</label>
          <textarea
            id="system-prompt"
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            rows={8}
            style={{ resize: 'vertical', fontFamily: 'inherit' }}
          />
        </div>

        <div className="actions">
          <button
            type="button"
            className="btn btn-primary"
            onClick={() => saveMutation.mutate()}
            disabled={saveMutation.isPending || !prompt.trim()}
          >
            {saveMutation.isPending ? 'Saving...' : 'Save prompt'}
          </button>
          <button
            type="button"
            className="btn btn-secondary"
            onClick={() => resetMutation.mutate()}
            disabled={resetMutation.isPending || !systemPrompt?.isCustom}
          >
            Reset to default
          </button>
        </div>

        {saveMutation.isError && (
          <div className="error" style={{ marginTop: '0.75rem' }}>
            Failed to save prompt
          </div>
        )}
        {saveMutation.isSuccess && (
          <div style={{ marginTop: '0.75rem', color: 'var(--success)' }}>Prompt saved</div>
        )}
      </div>
    </div>
  );
}
