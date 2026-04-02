import { useState, useEffect } from 'react';
import {
  FolderSearch, List, Copy, GitCompare, Trash2, Settings,
  Sun, Moon, Scan, Globe,
} from 'lucide-react';

import ScanView from './lib/views/ScanView';
import GlobalView from './lib/views/GlobalView';
import EntitiesView from './lib/views/EntitiesView';
import DiffView from './lib/views/DiffView';
import CopyView from './lib/views/CopyView';
import CleanView from './lib/views/CleanView';
import SettingsView from './lib/views/SettingsView';

type TabType = 'scan' | 'global' | 'entities' | 'diff' | 'copy' | 'clean' | 'settings';

const tabs = [
  { id: 'scan' as TabType, label: 'Scan', icon: FolderSearch },
  { id: 'global' as TabType, label: 'Global', icon: Globe },
  { id: 'entities' as TabType, label: 'Entities', icon: List },
  { id: 'diff' as TabType, label: 'Diff', icon: GitCompare },
  { id: 'copy' as TabType, label: 'Copy', icon: Copy },
  { id: 'clean' as TabType, label: 'Clean', icon: Trash2 },
  { id: 'settings' as TabType, label: 'Settings', icon: Settings },
] as const;

export default function App() {
  const [activeTab, setActiveTab] = useState<TabType>('scan');
  const [darkMode, setDarkMode] = useState(true);

  function toggleTheme() {
    const newDarkMode = !darkMode;
    setDarkMode(newDarkMode);
    if (newDarkMode) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }

  function handleNavigate(tab: string) {
    setActiveTab(tab as TabType);
  }

  useEffect(() => {
    if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
      setDarkMode(true);
      document.documentElement.classList.add('dark');
    }
  }, []);

  return (
    <div className="min-h-screen" style={{ backgroundColor: 'var(--bg-primary)', color: 'var(--text-primary)' }}>
      {/* Header */}
      <header className="sticky top-0 z-50 border-b" style={{ backgroundColor: 'var(--bg-primary)', borderColor: 'var(--border-color)' }}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ backgroundColor: 'var(--accent)' }}>
                <Scan size={22} className="text-white" />
              </div>
              <div>
                <h1 className="text-lg font-semibold" style={{ color: 'var(--text-primary)' }}>ccinspect</h1>
                <p className="text-xs" style={{ color: 'var(--text-muted)' }}>Claude Code Inspector</p>
              </div>
            </div>

            {/* Navigation Tabs */}
            <nav className="hidden md:flex items-center gap-1 p-1 rounded-xl" style={{ backgroundColor: 'var(--bg-muted)' }}>
              {tabs.map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className="px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200"
                  style={{
                    backgroundColor: activeTab === tab.id ? 'var(--bg-card)' : 'transparent',
                    color: activeTab === tab.id ? 'var(--accent)' : 'var(--text-secondary)',
                    boxShadow: activeTab === tab.id ? 'var(--shadow-card)' : 'none'
                  }}
                >
                  <div className="flex items-center gap-2">
                    <tab.icon size={16} />
                    {tab.label}
                  </div>
                </button>
              ))}
            </nav>

            {/* Actions */}
            <div className="flex items-center gap-2">
              <button
                onClick={toggleTheme}
                className="btn-icon"
                aria-label="Toggle theme"
              >
                {darkMode ? <Sun size={18} /> : <Moon size={18} />}
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {activeTab === 'scan' && <ScanView onNavigate={handleNavigate} />}
        {activeTab === 'global' && <GlobalView onNavigate={handleNavigate} />}
        {activeTab === 'entities' && <EntitiesView />}
        {activeTab === 'diff' && <DiffView />}
        {activeTab === 'copy' && <CopyView />}
        {activeTab === 'clean' && <CleanView />}
        {activeTab === 'settings' && <SettingsView />}
      </main>

      {/* Footer */}
      <footer className="border-t mt-12 py-6" style={{ borderColor: 'var(--border-color)' }}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex flex-col md:flex-row items-center justify-between gap-4">
            <p className="text-sm" style={{ color: 'var(--text-muted)' }}>
              Built with care by <span className="font-medium" style={{ color: 'var(--accent)' }}>@thieunv</span>
            </p>
            <div className="flex items-center gap-4">
              <span className="text-xs font-mono" style={{ color: 'var(--text-muted)' }}>v1.0.0</span>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
