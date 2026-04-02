import { useState } from 'react';
import { FolderOpen } from 'lucide-react';
import FileSystemBrowser from './FileSystemBrowser';

interface PathInputProps {
  value?: string;
  placeholder?: string;
  label?: string;
  allowGlobal?: boolean;
  onChange?: (value: string) => void;
}

export default function PathInput({
  value = '',
  placeholder = '/path/to/project',
  label,
  allowGlobal = true,
  onChange
}: PathInputProps) {
  const [showBrowser, setShowBrowser] = useState(false);
  const [inputValue, setInputValue] = useState(value);

  function openBrowser() {
    setShowBrowser(true);
  }

  function handlePathSelect(path: string) {
    setInputValue(path);
    onChange?.(path);
    setShowBrowser(false);
  }

  function handleInputChange(e: React.ChangeEvent<HTMLInputElement>) {
    const newValue = e.target.value;
    setInputValue(newValue);
    onChange?.(newValue);
  }

  function handleGlobalClick() {
    setInputValue('global');
    onChange?.('global');
  }

  return (
    <>
      <div className="space-y-1">
        {label && (
          <label className="block text-sm font-medium" style={{ color: 'var(--text-secondary)' }}>
            {label}
          </label>
        )}
        <div className="flex items-center gap-2">
          <input
            type="text"
            value={inputValue}
            onChange={handleInputChange}
            placeholder={placeholder}
            className="input-field flex-1"
          />
          {allowGlobal && (
            <button
              type="button"
              onClick={handleGlobalClick}
              className="px-3 py-2 text-xs font-medium rounded-md transition-all"
              style={{
                backgroundColor: 'var(--bg-muted)',
                color: 'var(--text-secondary)',
                border: '1px solid var(--border-subtle)'
              }}
              title="Use global config"
            >
              global
            </button>
          )}
          <button
            type="button"
            onClick={openBrowser}
            className="btn-icon"
            title="Browse folder"
          >
            <FolderOpen size={16} />
          </button>
        </div>
      </div>

      {showBrowser && (
        <FileSystemBrowser
          selectFolder={true}
          onSelect={handlePathSelect}
          onClose={() => setShowBrowser(false)}
        />
      )}
    </>
  );
}
