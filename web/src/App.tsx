import { useEffect, useMemo, useState } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import { BarChart3, ChevronLeft, LayoutGrid, Moon, Sun } from 'lucide-react';
import axios from 'axios';
import { DirectoryView } from './components/DirectoryView';
import ZiweiChart from './components/ZiweiChart';
import InterpretationPanel from './components/InterpretationPanel';
import PalaceDetailView from './components/PalaceDetailView';
import './styles/design-system.css';

interface BirthRecord {
  id: string;
  name: string;
  year: number;
  month: number;
  day: number;
  hour: number;
  gender: 'male' | 'female';
  is_lunar: boolean;
  is_leap: boolean;
  is_dst: boolean;
}

interface PalaceData {
  branch: string;
  palace_gan?: string;
  stars: string[];
  assistant_stars?: string[];
  secondary_stars?: string[];
  natal_transforms?: Array<{
    star: string;
    transformation: string;
    display: string;
  }>;
  liu_nian_stars?: string[];
  liu_nian_transforms?: Array<{
    star: string;
    transformation: string;
    display: string;
  }>;
  liu_yue_stars?: string[];
  liu_yue_transforms?: Array<{
    star: string;
    transformation: string;
    display: string;
  }>;
  liu_ri_stars?: string[];
  liu_ri_transforms?: Array<{
    star: string;
    transformation: string;
    display: string;
  }>;
  da_yun_ages?: string[];
  fly_hua?: {
    from_palace: string;
    stem: string;
    stages: Array<{
      type: string;
      star: string;
      target: string;
      motive?: string;
      action?: string;
      trap?: string;
    }>;
  };
}

interface KarmicStep {
  type: string;
  role: string;
  star: string;
  palace: string;
  desc: string;
}

interface SanFangRole {
  role: string;
  palace: string;
  diagnosis: string;
}

interface DeepStarAnalysis {
  name: string;
  verse: string;
  positive: string;
  negative: string;
  remedy: string;
  evolution?: string;
  brightness?: string;
}

interface FlyHuaAnalysis {
  from_palace: string;
  stem: string;
  stages: Array<{
    type: string;
    star: string;
    target: string;
    motive: string;
    action: string;
    trap: string;
  }>;
}

interface InterpretationData {
  summary: string;
  character_traits: string;
  origin_palace_analysis: string;
  karmic_narrative: KarmicStep[];
  san_fang_diagnosis: SanFangRole[];
  star_details?: DeepStarAnalysis[];
  origin_fly_hua?: FlyHuaAnalysis;
  temporal_resonance: Array<{
    layer: string;
    type: string;
    star: string;
    natal?: string;
    palace: string;
    mood: string;
  }>;
  classic_patterns?: string[];
}

interface ZiweiResponse {
  gender: string;
  wuxing: string;
  na_yin: string;
  origin_palace: string;
  ming_gong: string;
  shen_gong: string;
  year_pillar: string;
  day_pillar: string;
  current_da_yun?: {
    index: number;
    start_age: number;
    end_age: number;
    stem: string;
    branch: string;
    palace: string;
  } | null;
  da_yun?: Array<{
    index: number;
    start_age: number;
    end_age: number;
    stem: string;
    branch: string;
    palace: string;
  }>;
  liu_nian?: {
    label: string;
    branch: string;
    palace: string;
    stem: string;
  } | null;
  liu_yue?: {
    label: string;
    branch: string;
    palace: string;
    stem: string;
  } | null;
  liu_ri?: {
    label: string;
    branch: string;
    palace: string;
    stem: string;
  } | null;
  palaces: Record<string, PalaceData>;
  interpretation: InterpretationData;
}

type ResultTab = 'overview' | 'chart';

const resultTabs: Array<{ key: ResultTab; label: string; icon: typeof LayoutGrid }> = [
  { key: 'overview', label: '總覽', icon: LayoutGrid },
  { key: 'chart', label: '命盤', icon: BarChart3 },
];

type FocusedPalace = {
  label: string;
  palaceName: string;
  branch?: string;
} | null;

const ResultsPage = ({
  data,
  onBack,
  darkMode,
  setDarkMode,
  userName,
}: {
  data: ZiweiResponse;
  onBack: () => void;
  darkMode: boolean;
  setDarkMode: (dark: boolean) => void;
  userName: string;
}) => {
  const [activeTab, setActiveTab] = useState<ResultTab>('overview');
  const [focusedPalace, setFocusedPalace] = useState<FocusedPalace>(null);

  const palaceByBranch = useMemo(() => {
    const map = new Map<string, string>();
    Object.entries(data.palaces).forEach(([palaceName, palace]) => {
      if (palace.branch) {
        map.set(palace.branch, palaceName);
      }
    });
    return map;
  }, [data.palaces]);

  const branchByPalace = useMemo(() => {
    const map = new Map<string, string>();
    Object.entries(data.palaces).forEach(([palaceName, palace]) => {
      if (palace.branch) {
        map.set(palaceName, palace.branch);
      }
    });
    return map;
  }, [data.palaces]);

  const keyInsights = useMemo(() => {
    const mingHost = palaceByBranch.get(data.ming_gong);
    const shenHost = palaceByBranch.get(data.shen_gong);
    const originBranch = branchByPalace.get(data.origin_palace);

    return [
      {
        label: `命宮（${data.ming_gong}）`,
        value: mingHost || data.ming_gong,
        palace: mingHost,
        branch: data.ming_gong,
      },
      {
        label: `身宮（${data.shen_gong}）`,
        value: shenHost || data.shen_gong,
        palace: shenHost,
        branch: data.shen_gong,
      },
      {
        label: '五行局',
        value: data.wuxing,
      },
      {
        label: `來因宮（${originBranch || ''}）`,
        value: data.origin_palace,
        palace: data.origin_palace,
        branch: originBranch,
      },
      {
        label: '納音',
        value: data.na_yin,
      },
      {
        label: '年柱',
        value: data.year_pillar,
      },
    ];
  }, [data, palaceByBranch, branchByPalace]);

  return (
    <div className="page-shell">
      <div className="page-container section-stack">
        <motion.header initial={{ opacity: 0, y: -16 }} animate={{ opacity: 1, y: 0 }} className="card section-stack">
          <div style={{ display: 'flex', justifyContent: 'space-between', gap: '1rem', flexWrap: 'wrap', alignItems: 'center' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem', flexWrap: 'wrap' }}>
              <button className="btn-secondary" onClick={onBack}>
                <ChevronLeft size={16} />
                返回資料列表
              </button>
              <div>
                <div className="heading-lg">{userName} 命盤</div>
                <div className="body-sm">先看總覽，再看命盤，最後深入解讀。</div>
              </div>
            </div>
            <button className="btn-secondary" onClick={() => setDarkMode(!darkMode)}>
              {darkMode ? <Sun size={16} /> : <Moon size={16} />}
              {darkMode ? '淺色' : '深色'}
            </button>
          </div>

          <div className="metric-grid">
            {keyInsights.map((item) => (
              <button
                key={item.label}
                className="metric-card"
                style={{ 
                  cursor: item.palace ? 'pointer' : 'default',
                  textAlign: 'left',
                  border: 'none',
                  background: 'var(--surface)',
                }}
                onClick={() => {
                  if (item.palace) {
                    setFocusedPalace({
                      label: item.label,
                      palaceName: item.palace,
                      branch: item.branch,
                    });
                  }
                }}
                disabled={!item.palace}
                title={item.palace ? `點擊查看 ${item.label} 的三方四正` : undefined}
              >
                <div className="metric-label">{item.label}</div>
                <div className="metric-value">{item.value}</div>
              </button>
            ))}
          </div>

          <div className="card-gold card" style={{ padding: '1rem 1.25rem' }}>
            <div className="heading-sm" style={{ marginBottom: '0.5rem' }}>總結導讀</div>
            <div className="body-md">{data.interpretation.summary}</div>
          </div>
        </motion.header>

        <section className="section-stack">
          {focusedPalace ? (
            <PalaceDetailView
              palaces={data.palaces}
              focus={focusedPalace}
              onBack={() => setFocusedPalace(null)}
            />
          ) : (
            <>
              <div className="tab-row">
                {resultTabs.map((tab) => {
                  const Icon = tab.icon;
                  return (
                    <button key={tab.key} className={`tab-button ${activeTab === tab.key ? 'is-active' : ''}`} onClick={() => setActiveTab(tab.key)}>
                      <Icon size={14} />
                      {tab.label}
                    </button>
                  );
                })}
              </div>

              {activeTab === 'overview' && (
                <motion.div initial={{ opacity: 0, y: 12 }} animate={{ opacity: 1, y: 0 }}>
                  <InterpretationPanel interpretation={data.interpretation} />
                </motion.div>
              )}

              {activeTab === 'chart' && (
                <motion.div initial={{ opacity: 0, y: 12 }} animate={{ opacity: 1, y: 0 }}>
                  <ZiweiChart
                    palaces={data.palaces}
                    mingGong={data.ming_gong}
                    shenGong={data.shen_gong}
                    originPalace={data.origin_palace}
                    currentDaYun={data.current_da_yun}
                    daYun={data.da_yun}
                    liuNian={data.liu_nian}
                    liuYue={data.liu_yue}
                    liuRi={data.liu_ri}
                  />
                </motion.div>
              )}
            </>
          )}
        </section>
      </div>
    </div>
  );
};

export default function App() {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<ZiweiResponse | null>(null);
  const [darkMode, setDarkMode] = useState(false);
  const [currentPage, setCurrentPage] = useState<'directory' | 'results'>('directory');
  const [currentUser, setCurrentUser] = useState<BirthRecord | null>(null);

  useEffect(() => {
    if (darkMode) {
      document.body.classList.add('dark-theme');
      return;
    }
    document.body.classList.remove('dark-theme');
  }, [darkMode]);

  const calculate = async (record: BirthRecord) => {
    setLoading(true);
    setCurrentUser(record);
    try {
      const payload = {
        year: record.year,
        month: record.month,
        day: record.day,
        hour: record.hour,
        gender: record.gender,
        is_lunar: record.is_lunar,
        is_leap: record.is_leap,
        is_dst: record.is_dst,
      };
      const resp = await axios.post('/api/v1/calculate', payload);
      setData(resp.data);
      setCurrentPage('results');
    } catch (err) {
      console.error(err);
      alert('計算失敗');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={darkMode ? 'dark-theme' : ''}>
      <AnimatePresence mode="wait">
        {currentPage === 'directory' && (
          <motion.div key="directory" initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}>
            <DirectoryView onSelect={calculate} />
          </motion.div>
        )}
        {currentPage === 'results' && data && (
          <motion.div key="results" initial={{ opacity: 0, y: 12 }} animate={{ opacity: 1, y: 0 }} exit={{ opacity: 0 }}>
            <ResultsPage
              data={data}
              onBack={() => setCurrentPage('directory')}
              darkMode={darkMode}
              setDarkMode={setDarkMode}
              userName={currentUser?.name || '未知命主'}
            />
          </motion.div>
        )}
      </AnimatePresence>

      {loading && (
        <div style={{ position: 'fixed', inset: 0, background: 'rgba(24, 18, 14, 0.38)', backdropFilter: 'blur(6px)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 120 }}>
          <div className="card section-stack" style={{ alignItems: 'center', minWidth: '260px' }}>
            <div style={{ width: '3rem', height: '3rem', border: '4px solid rgba(178, 135, 70, 0.2)', borderTopColor: 'var(--cta)', borderRadius: '999px', animation: 'spin 1s linear infinite' }} />
            <div className="heading-sm">正在排盤</div>
            <div className="body-sm">系統正在計算命盤與分析內容。</div>
          </div>
        </div>
      )}
    </div>
  );
}
