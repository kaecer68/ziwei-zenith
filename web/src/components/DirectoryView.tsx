import { useEffect, useMemo, useState } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import {
  ArrowRight,
  Clock,
  HelpCircle,
  History as HistoryIcon,
  Search,
  ShieldCheck,
  Trash2,
  X,
} from 'lucide-react';
import axios from 'axios';

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
  tags: string[];
  notes?: string;
  created_at: string;
}

interface TagType {
  id: string;
  name: string;
  color: string;
}

interface DirectoryViewProps {
  onSelect: (record: BirthRecord) => void;
}

type StepKey = 'identity' | 'birth' | 'advanced' | 'confirm';

const stepOrder: Array<{ key: StepKey; title: string; desc: string }> = [
  { key: 'identity', title: '基本身份', desc: '輸入姓名與性別' },
  { key: 'birth', title: '出生時間', desc: '設定曆法、日期與時辰' },
  { key: 'advanced', title: '進階校準', desc: '處理 DST 與農曆閏月' },
  { key: 'confirm', title: '確認排盤', desc: '確認資料後建立命盤' },
];

export const DirectoryView = ({ onSelect }: DirectoryViewProps) => {
  const [records, setRecords] = useState<BirthRecord[]>([]);
  const [tags, setTags] = useState<TagType[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [selectedTag, setSelectedTag] = useState<string | null>(null);
  const [isDstModalOpen, setIsDstModalOpen] = useState(false);
  const [activeStep, setActiveStep] = useState<StepKey>('identity');
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [quickAdd, setQuickAdd] = useState({
    name: '',
    year: new Date().getFullYear() - 30,
    month: 1,
    day: 1,
    hour: 12,
    gender: 'male' as 'male' | 'female',
    is_lunar: false,
    is_leap: false,
    is_dst: false,
  });

  const fetchData = async () => {
    setLoading(true);
    try {
      const [recRes, tagRes] = await Promise.all([
        axios.get('http://localhost:8081/api/v1/records'),
        axios.get('http://localhost:8081/api/v1/tags'),
      ]);
      setRecords(recRes.data || []);
      setTags(tagRes.data || []);
    } catch (err) {
      console.error('Failed to fetch directory data', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const filteredRecords = useMemo(() => {
    return records.filter((record) => {
      const nameMatch = (record.name || '').toLowerCase().includes(search.toLowerCase());
      return nameMatch && (!selectedTag || (record.tags || []).includes(selectedTag));
    });
  }, [records, search, selectedTag]);

  const stepIndex = stepOrder.findIndex((step) => step.key === activeStep);
  const summaryLabel = `${quickAdd.is_lunar ? '農曆' : '公曆'} ${quickAdd.year}年 ${quickAdd.month}月 ${quickAdd.day}日 ${quickAdd.hour}時`;

  const handleQuickAdd = async () => {
    if (!quickAdd.name.trim()) {
      alert('請先輸入姓名');
      setActiveStep('identity');
      return;
    }
    try {
      const payload = {
        ...quickAdd,
        tags: selectedTags,
        notes: '',
      };
      const response = await axios.post('http://localhost:8081/api/v1/records', payload);
      const newRecord = {
        ...payload,
        id: response.data.id || `temp-${Date.now()}`,
        created_at: new Date().toISOString(),
      } as BirthRecord;
      onSelect(newRecord);
      setQuickAdd((prev) => ({ ...prev, name: '' }));
      setSelectedTags([]);
      fetchData();
    } catch (err) {
      console.error(err);
      alert('保存失敗，請檢查系統連接');
    }
  };

  const handleDelete = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    if (!confirm('此操作無法撤銷。確定刪除紀錄？')) return;
    try {
      await axios.delete(`http://localhost:8081/api/v1/records/${id}`);
      fetchData();
    } catch (err) {
      console.error(err);
    }
  };

  const renderStep = () => {
    if (activeStep === 'identity') {
      return (
        <div className="section-stack">
          <div>
            <div className="heading-md">先確認命主身份</div>
            <div className="body-sm">先填最基本的資料，讓整個排盤流程更聚焦。</div>
          </div>
          <div className="form-grid">
            <div className="field">
              <label htmlFor="name">姓名</label>
              <input
                id="name"
                value={quickAdd.name}
                onChange={(e) => setQuickAdd({ ...quickAdd, name: e.target.value })}
                placeholder="例如：王小明"
              />
              <small>用於辨識命盤紀錄，可稍後再編輯。</small>
            </div>
            <div className="field">
              <label>性別</label>
              <div className="toggle-row">
                <button className={`toggle-chip ${quickAdd.gender === 'male' ? 'is-active' : ''}`} onClick={() => setQuickAdd({ ...quickAdd, gender: 'male' })}>
                  乾造（男命）
                </button>
                <button className={`toggle-chip ${quickAdd.gender === 'female' ? 'is-active' : ''}`} onClick={() => setQuickAdd({ ...quickAdd, gender: 'female' })}>
                  坤造（女命）
                </button>
              </div>
              <small>影響部分運限計算方向。</small>
            </div>
            {tags.length > 0 && (
              <div className="field">
                <label>分類標籤</label>
                <div className="toggle-row">
                  {tags.map((tag) => {
                    const isSelected = selectedTags.includes(tag.name);
                    return (
                      <button
                        key={tag.id}
                        className={`toggle-chip ${isSelected ? 'is-active' : ''}`}
                        style={isSelected ? { borderColor: tag.color, color: tag.color } : {}}
                        onClick={() =>
                          setSelectedTags((prev) =>
                            isSelected ? prev.filter((t) => t !== tag.name) : [...prev, tag.name]
                          )
                        }
                      >
                        {tag.name}
                      </button>
                    );
                  })}
                </div>
                <small>可多選，方便後續在紀錄庫快速篩選。</small>
              </div>
            )}
          </div>
        </div>
      );
    }

    if (activeStep === 'birth') {
      return (
        <div className="section-stack">
          <div>
            <div className="heading-md">輸入出生時間</div>
            <div className="body-sm">先切換公曆或農曆，再輸入年月日與時辰。</div>
          </div>
          <div className="field">
            <label>曆法</label>
            <div className="toggle-row">
              <button className={`toggle-chip ${!quickAdd.is_lunar ? 'is-active' : ''}`} onClick={() => setQuickAdd({ ...quickAdd, is_lunar: false, is_leap: false })}>
                公曆
              </button>
              <button className={`toggle-chip ${quickAdd.is_lunar ? 'is-active' : ''}`} onClick={() => setQuickAdd({ ...quickAdd, is_lunar: true })}>
                農曆
              </button>
            </div>
          </div>
          <div className="form-grid">
            <div className="field">
              <label htmlFor="year">年份</label>
              <input id="year" type="number" value={quickAdd.year} onChange={(e) => setQuickAdd({ ...quickAdd, year: Number(e.target.value) })} />
            </div>
            <div className="field">
              <label htmlFor="month">月份</label>
              <input id="month" type="number" min="1" max="12" value={quickAdd.month} onChange={(e) => setQuickAdd({ ...quickAdd, month: Number(e.target.value) })} />
            </div>
            <div className="field">
              <label htmlFor="day">日期</label>
              <input id="day" type="number" min="1" max="31" value={quickAdd.day} onChange={(e) => setQuickAdd({ ...quickAdd, day: Number(e.target.value) })} />
            </div>
            <div className="field">
              <label htmlFor="hour">時辰（24 小時）</label>
              <input id="hour" type="number" min="0" max="23" value={quickAdd.hour} onChange={(e) => setQuickAdd({ ...quickAdd, hour: Number(e.target.value) })} />
            </div>
          </div>
        </div>
      );
    }

    if (activeStep === 'advanced') {
      return (
        <div className="section-stack">
          <div>
            <div className="heading-md">進階校準</div>
            <div className="body-sm">只有在特殊情況下才需要調整，不懂可先維持預設。</div>
          </div>
          <div className="card">
            <div style={{ display: 'flex', justifyContent: 'space-between', gap: '1rem', alignItems: 'center', flexWrap: 'wrap' }}>
              <div>
                <div className="heading-sm">夏令時間</div>
                <div className="body-sm">若出生地當時有實施夏令時間，需勾選後進行校正。</div>
              </div>
              <div className="toggle-row">
                <button className={`toggle-chip ${quickAdd.is_dst ? 'is-active' : ''}`} onClick={() => setQuickAdd({ ...quickAdd, is_dst: !quickAdd.is_dst })}>
                  {quickAdd.is_dst ? '已啟用 DST' : '未啟用 DST'}
                </button>
                <button className="btn-secondary" onClick={() => setIsDstModalOpen(true)}>
                  <HelpCircle size={16} />
                  查看說明
                </button>
              </div>
            </div>
          </div>
          {quickAdd.is_lunar && (
            <div className="card">
              <div style={{ display: 'flex', justifyContent: 'space-between', gap: '1rem', alignItems: 'center', flexWrap: 'wrap' }}>
                <div>
                  <div className="heading-sm">農曆閏月</div>
                  <div className="body-sm">只有在農曆出生資料確定為閏月時才需要啟用。</div>
                </div>
                <button className={`toggle-chip ${quickAdd.is_leap ? 'is-active' : ''}`} onClick={() => setQuickAdd({ ...quickAdd, is_leap: !quickAdd.is_leap })}>
                  {quickAdd.is_leap ? '已標記閏月' : '不是閏月'}
                </button>
              </div>
            </div>
          )}
        </div>
      );
    }

    return (
      <div className="section-stack">
        <div>
          <div className="heading-md">確認後開始排盤</div>
          <div className="body-sm">最後再確認一次資料，避免命盤偏差。</div>
        </div>
        <div className="card-gold card">
          <div className="section-stack">
            <div className="metric-grid" style={{ gridTemplateColumns: 'repeat(2, minmax(0, 1fr))' }}>
              <div className="metric-card">
                <div className="metric-label">命主姓名</div>
                <div className="metric-value">{quickAdd.name || '未填寫'}</div>
              </div>
              <div className="metric-card">
                <div className="metric-label">性別</div>
                <div className="metric-value">{quickAdd.gender === 'male' ? '乾造' : '坤造'}</div>
              </div>
              <div className="metric-card">
                <div className="metric-label">出生資料</div>
                <div className="metric-value" style={{ fontSize: '1.05rem' }}>{summaryLabel}</div>
              </div>
              <div className="metric-card">
                <div className="metric-label">進階校準</div>
                <div className="metric-value" style={{ fontSize: '1.05rem' }}>
                  {quickAdd.is_dst ? 'DST 已開啟' : 'DST 關閉'}{quickAdd.is_lunar ? quickAdd.is_leap ? ' / 閏月' : ' / 非閏月' : ''}
                </div>
              </div>
            </div>
            <div className="body-sm">你可以回前一步調整，或直接建立命盤。</div>
            <div style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap' }}>
              <button className="btn-secondary" onClick={() => setActiveStep('birth')}>返回修改</button>
              <button className="btn-primary" onClick={handleQuickAdd}>
                開始排盤
                <ArrowRight size={16} />
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className="page-shell">
      <motion.div className="page-container section-stack" initial={{ opacity: 0, y: 12 }} animate={{ opacity: 1, y: 0 }}>
        <section className="hero-panel">
          <div className="card section-stack">
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.65rem' }}>
              <ShieldCheck size={18} color="var(--cta)" />
              <span className="body-sm">紫微斗數排盤與分析工具</span>
            </div>
            <div className="heading-xl">紫微大觀</div>
            <div className="body-lg">
              把建檔、排盤與解讀整理成一條清楚流程，讓第一次接觸紫微斗數的人也能順利完成命盤建立。
            </div>
          </div>
          <div className="card-gold card section-stack">
            <div className="heading-md">這次你可以做什麼</div>
            <div className="body-sm">先建立一筆新命盤，或從下方紀錄庫快速回看既有命盤。</div>
            <div className="metric-grid" style={{ gridTemplateColumns: 'repeat(2, minmax(0, 1fr))' }}>
              <div className="metric-card">
                <div className="metric-label">已存命盤</div>
                <div className="metric-value">{records.length}</div>
              </div>
              <div className="metric-card">
                <div className="metric-label">標籤數量</div>
                <div className="metric-value">{tags.length}</div>
              </div>
            </div>
          </div>
        </section>

        <section className="step-shell">
          <div className="step-nav">
            {stepOrder.map((step, index) => (
              <button key={step.key} className={`step-nav-item ${step.key === activeStep ? 'is-active' : ''}`} onClick={() => setActiveStep(step.key)}>
                <span className="step-index">{index + 1}</span>
                <span>
                  <div className="heading-sm">{step.title}</div>
                  <div className="body-sm">{step.desc}</div>
                </span>
              </button>
            ))}
          </div>

          <section className="step-panel section-stack">
            {renderStep()}
            <div style={{ display: 'flex', justifyContent: 'space-between', gap: '0.75rem', flexWrap: 'wrap' }}>
              <button className="btn-secondary" disabled={stepIndex === 0} onClick={() => setActiveStep(stepOrder[Math.max(0, stepIndex - 1)].key)}>
                上一步
              </button>
              {activeStep !== 'confirm' && (
                <button className="btn-primary" onClick={() => setActiveStep(stepOrder[Math.min(stepOrder.length - 1, stepIndex + 1)].key)}>
                  下一步
                  <ArrowRight size={16} />
                </button>
              )}
            </div>
          </section>
        </section>

        <section className="section-stack">
          <div className="archive-toolbar">
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
              <HistoryIcon size={18} color="var(--cta)" />
              <div>
                <div className="heading-md">命盤紀錄庫</div>
                <div className="body-sm">快速搜尋、回看與管理已建立的命盤。</div>
              </div>
            </div>
            <div style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap' }}>
              <div className="field" style={{ minWidth: '260px' }}>
                <label htmlFor="search">搜尋姓名</label>
                <div style={{ position: 'relative' }}>
                  <Search size={16} style={{ position: 'absolute', left: '0.85rem', top: '50%', transform: 'translateY(-50%)', color: 'var(--text-soft)' }} />
                  <input id="search" style={{ paddingLeft: '2.35rem' }} value={search} onChange={(e) => setSearch(e.target.value)} placeholder="輸入姓名查找紀錄" />
                </div>
              </div>
            </div>
          </div>

          {tags.length > 0 && (
            <div className="tab-row">
              <button className={`tab-button ${selectedTag === null ? 'is-active' : ''}`} onClick={() => setSelectedTag(null)}>
                全部
              </button>
              {tags.map((tag) => (
                <button key={tag.id} className={`tab-button ${selectedTag === tag.name ? 'is-active' : ''}`} onClick={() => setSelectedTag(selectedTag === tag.name ? null : tag.name)}>
                  {tag.name}
                </button>
              ))}
            </div>
          )}

          {loading ? (
            <div className="empty-state">正在載入命盤紀錄…</div>
          ) : filteredRecords.length === 0 ? (
            <div className="empty-state">目前沒有符合條件的命盤紀錄。</div>
          ) : (
            <div className="archive-grid">
              <AnimatePresence>
                {filteredRecords.map((record) => (
                  <motion.div key={record.id} layout initial={{ opacity: 0, y: 16 }} animate={{ opacity: 1, y: 0 }} exit={{ opacity: 0, scale: 0.98 }} className="archive-card" onClick={() => onSelect(record)}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', gap: '1rem' }}>
                      <div className="section-stack" style={{ gap: '0.75rem' }}>
                        <div>
                          <div className="heading-md" style={{ fontSize: '1.25rem' }}>{record.name}</div>
                          <div className="body-sm">{record.gender === 'male' ? '乾造' : '坤造'} · {record.is_lunar ? '農曆' : '公曆'}</div>
                        </div>
                        <div className="body-sm" style={{ display: 'flex', alignItems: 'center', gap: '0.45rem' }}>
                          <Clock size={14} />
                          {record.year}年 {record.month}月 {record.day}日 {record.hour}時
                        </div>
                        <div className="tab-row">
                          {record.is_dst && <span className="tab-button is-active">DST</span>}
                          {record.tags?.map((tag) => (
                            <span key={tag} className="tab-button">{tag}</span>
                          ))}
                        </div>
                      </div>
                      <button className="btn-secondary" onClick={(e) => handleDelete(record.id, e)}>
                        <Trash2 size={16} />
                      </button>
                    </div>
                  </motion.div>
                ))}
              </AnimatePresence>
            </div>
          )}
        </section>
      </motion.div>

      <AnimatePresence>
        {isDstModalOpen && (
          <div style={{ position: 'fixed', inset: 0, background: 'rgba(24, 18, 14, 0.45)', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '1rem', zIndex: 100 }}>
            <motion.div initial={{ opacity: 0, y: 12 }} animate={{ opacity: 1, y: 0 }} exit={{ opacity: 0, y: 12 }} className="card" style={{ width: 'min(760px, 100%)', position: 'relative' }}>
              <button className="btn-secondary" style={{ position: 'absolute', top: '1rem', right: '1rem' }} onClick={() => setIsDstModalOpen(false)}>
                <X size={16} />
              </button>
              <div className="section-stack">
                <div className="heading-md">夏令時間校準說明</div>
                <div className="body-md">
                  若出生地在該時期有實施夏令時間，排盤前要先扣回 1 小時。若不確定，可先查詢當地歷史制度，再決定是否勾選。
                </div>
                <div className="insight-grid">
                  <div className="card-gold card">
                    <div className="heading-sm">台灣</div>
                    <div className="body-sm">1945~1961、1974~1975、1979 曾有夏令時間實施紀錄。</div>
                  </div>
                  <div className="card-gold card">
                    <div className="heading-sm">中國大陸</div>
                    <div className="body-sm">1986~1991 固定實施，通常落在每年 4~9 月之間。</div>
                  </div>
                </div>
                <div className="card" style={{ display: 'flex', gap: '0.75rem', alignItems: 'flex-start' }}>
                  <HelpCircle size={18} color="var(--cta)" />
                  <div className="body-sm">若出生時間接近子時，DST 校正可能影響日柱判斷，建議務必確認。</div>
                </div>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>
    </div>
  );
};
