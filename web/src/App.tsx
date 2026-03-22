import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
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
    time_branch: string;
  } | null;
  liu_yue?: {
    label: string;
    branch: string;
    palace: string;
    stem: string;
    time_branch: string;
  } | null;
  liu_ri?: {
    label: string;
    branch: string;
    palace: string;
    stem: string;
    time_branch: string;
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
  temporalInfo?: {
    stem: string;
    timeBranch?: string;
    ageRange?: string;
    timeLabel?: string;
  };
} | null;

type TemporalPositionSummary = {
  key: 'da_yun' | 'liu_nian' | 'liu_yue' | 'liu_ri';
  label: string;
  palace?: string;
  branch?: string;
  timeGanzhi?: string;
  detail?: string;
};

const ResultsPage = ({
  data,
  onBack,
  darkMode,
  setDarkMode,
  userInfo,
}: {
  data: ZiweiResponse;
  onBack: () => void;
  darkMode: boolean;
  setDarkMode: (dark: boolean) => void;
  userInfo: BirthRecord | null;
}) => {
  const [activeTab, setActiveTab] = useState<ResultTab>('overview');
  const [focusedPalace, setFocusedPalace] = useState<FocusedPalace>(null);
  const [selectedInsight, setSelectedInsight] = useState<string | null>(null);
  
  // 大限切換狀態
  const [currentDaYunIndex, setCurrentDaYunIndex] = useState<number>(0);
  const [currentTemporalData, setCurrentTemporalData] = useState<{
    da_yun: typeof data.current_da_yun;
    liu_nian: typeof data.liu_nian;
    liu_yue: typeof data.liu_yue;
    liu_ri: typeof data.liu_ri;
  } | null>(null);
  const [targetDate, setTargetDate] = useState<Date>(() => new Date());
  const [isTemporalLoading, setIsTemporalLoading] = useState(false);
  const temporalRequestSeq = useRef(0);
  const targetYear = targetDate.getFullYear();
  const targetMonth = targetDate.getMonth() + 1;
  const targetDay = targetDate.getDate();

  const getEffectiveFlowYear = useCallback((date: Date) => {
    const month = date.getMonth() + 1;
    const day = date.getDate();
    return month < 2 || (month === 2 && day < 4) ? date.getFullYear() - 1 : date.getFullYear();
  }, []);

  const effectiveFlowYear = useMemo(() => getEffectiveFlowYear(targetDate), [targetDate, getEffectiveFlowYear]);

  const findDaYunIndexByAge = useCallback((age: number) => {
    if (!data.da_yun || data.da_yun.length === 0) return 0;
    for (let i = 0; i < data.da_yun.length; i++) {
      const dy = data.da_yun[i];
      if (age >= dy.start_age && age <= dy.end_age) {
        return i;
      }
    }
    return 0;
  }, [data.da_yun]);

  const syncDaYunWithDate = useCallback((date: Date) => {
    if (!userInfo || !data.da_yun || data.da_yun.length === 0) return;
    const flowYear = getEffectiveFlowYear(date);
    const virtualAge = Math.max(1, flowYear - userInfo.year + 1);
    const idx = findDaYunIndexByAge(virtualAge);
    if (idx !== currentDaYunIndex) {
      setCurrentDaYunIndex(idx);
    }
  }, [userInfo, data.da_yun, findDaYunIndexByAge, currentDaYunIndex, getEffectiveFlowYear]);

  useEffect(() => {
    syncDaYunWithDate(targetDate);
  }, [targetDate, syncDaYunWithDate]);

  // 獲取指定大限的運限數據
  const fetchTemporalData = useCallback(async (daYunIndex: number, year: number, month: number, day: number) => {
    if (!userInfo) return;

    const requestSeq = ++temporalRequestSeq.current;
    setIsTemporalLoading(true);
    try {
      const response = await axios.post('/api/v1/calculate/temporal', {
        birth_year: userInfo.year,
        birth_month: userInfo.month,
        birth_day: userInfo.day,
        birth_hour: userInfo.hour,
        gender: userInfo.gender,
        is_lunar: userInfo.is_lunar,
        is_leap: userInfo.is_leap,
        da_yun_index: daYunIndex,
        target_year: year,
        target_month: month,
        target_day: day,
      });

      if (requestSeq === temporalRequestSeq.current) {
        setCurrentTemporalData(response.data);
      }
    } catch (err) {
      console.error('獲取運限數據失敗:', err);
    } finally {
      if (requestSeq === temporalRequestSeq.current) {
        setIsTemporalLoading(false);
      }
    }
  }, [userInfo]);

  // 當大限索引改變時，獲取新的運限數據
  useEffect(() => {
    fetchTemporalData(currentDaYunIndex, targetYear, targetMonth, targetDay);
  }, [currentDaYunIndex, targetYear, targetMonth, targetDay, fetchTemporalData]);

  // 切換大限的函數
  const handlePrevDaYun = () => {
    if (currentDaYunIndex > 0) {
      setCurrentDaYunIndex(currentDaYunIndex - 1);
    }
  };

  const handleNextDaYun = () => {
    if (data.da_yun && currentDaYunIndex < data.da_yun.length - 1) {
      setCurrentDaYunIndex(currentDaYunIndex + 1);
    }
  };

  const handleDaYunClick = (index: number) => {
    setCurrentDaYunIndex(index);
    if (!userInfo || !data.da_yun || !data.da_yun[index]) return;

    const selectedDaYun = data.da_yun[index];
    const startYear = userInfo.year + selectedDaYun.start_age - 1;
    const endYear = userInfo.year + selectedDaYun.end_age - 1;

    setTargetDate((prev) => {
      const flowYear = getEffectiveFlowYear(prev);
      if (flowYear >= startYear && flowYear <= endYear) return prev;
      const next = new Date(prev);
      next.setFullYear(startYear);
      return next;
    });
  };

  const handleFlowYearSelect = (year: number) => {
    setTargetDate((prev) => {
      const next = new Date(prev);
      next.setFullYear(year);
      return next;
    });
  };

  const handleFlowYearShift = (delta: number) => {
    setTargetDate((prev) => {
      const next = new Date(prev);
      next.setFullYear(prev.getFullYear() + delta);
      return next;
    });
  };

  const handleFlowMonthSelect = (month: number) => {
    setTargetDate((prev) => {
      const currentDay = prev.getDate();
      const next = new Date(prev);
      next.setDate(1);
      next.setMonth(month - 1);
      const maxDay = new Date(next.getFullYear(), month, 0).getDate();
      next.setDate(Math.min(currentDay, maxDay));
      return next;
    });
  };

  const handleFlowDaySelect = (selectedDate: Date) => {
    setTargetDate(new Date(selectedDate));
  };

  const handleFlowDayShift = (delta: number) => {
    setTargetDate((prev) => {
      const next = new Date(prev);
      next.setDate(prev.getDate() + delta * 12);
      return next;
    });
  };

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

  const selectedDaYun = useMemo(() => {
    if (!data.da_yun || !data.da_yun[currentDaYunIndex]) return null;
    return data.da_yun[currentDaYunIndex];
  }, [data.da_yun, currentDaYunIndex]);

  const displayDaYun = currentTemporalData?.da_yun || selectedDaYun || data.current_da_yun;
  const displayLiuNian = currentTemporalData?.liu_nian || data.liu_nian;
  const displayLiuYue = currentTemporalData?.liu_yue || data.liu_yue;
  const displayLiuRi = currentTemporalData?.liu_ri || data.liu_ri;

  const formatTemporalGanzhi = (temporal?: { stem?: string; time_branch?: string } | null) => {
    if (!temporal?.stem || !temporal?.time_branch) return '未知';
    return `${temporal.stem}${temporal.time_branch}`;
  };

  const flowYearOptions = useMemo(() => {
    if (!selectedDaYun || !userInfo) return [] as number[];
    return Array.from(
      { length: selectedDaYun.end_age - selectedDaYun.start_age + 1 },
      (_, i) => userInfo.year + selectedDaYun.start_age - 1 + i,
    );
  }, [selectedDaYun, userInfo]);

  const flowMonthOptions = useMemo(() => Array.from({ length: 12 }, (_, i) => i + 1), []);

  const flowDayOptions = useMemo(() => {
    return Array.from({ length: 12 }, (_, i) => {
      const dayOffset = i - 5;
      const optionDate = new Date(targetDate);
      optionDate.setDate(targetDate.getDate() + dayOffset);
      const isCurrentMonth =
        optionDate.getFullYear() === targetDate.getFullYear() &&
        optionDate.getMonth() === targetDate.getMonth();

      return {
        date: optionDate,
        day: optionDate.getDate(),
        month: optionDate.getMonth() + 1,
        isCurrentMonth,
      };
    });
  }, [targetDate]);

  const temporalLayers = [
    {
      label: '大限',
      palace: displayDaYun?.palace,
      branch: displayDaYun?.branch,
      stem: displayDaYun?.stem,
      timeBranch: displayDaYun?.branch,
      timeLabel: `${effectiveFlowYear}流年`,
    },
    {
      label: '流年',
      palace: displayLiuNian?.palace,
      branch: displayLiuNian?.branch,
      stem: displayLiuNian?.stem,
      timeBranch: displayLiuNian?.time_branch,
      timeLabel: `${effectiveFlowYear}年`,
    },
    {
      label: '流月',
      palace: displayLiuYue?.palace,
      branch: displayLiuYue?.branch,
      stem: displayLiuYue?.stem,
      timeBranch: displayLiuYue?.time_branch,
      timeLabel: `${targetMonth}月`,
    },
    {
      label: '流日',
      palace: displayLiuRi?.palace,
      branch: displayLiuRi?.branch,
      stem: displayLiuRi?.stem,
      timeBranch: displayLiuRi?.time_branch,
      timeLabel: `${targetMonth}月${targetDay}日`,
    },
  ];

  const temporalPositions = useMemo<TemporalPositionSummary[]>(() => {
    const list: TemporalPositionSummary[] = [];

    if (displayDaYun) {
      list.push({
        key: 'da_yun',
        label: '大限',
        palace: displayDaYun.palace,
        branch: displayDaYun.branch,
        timeGanzhi: displayDaYun.stem && displayDaYun.branch ? `${displayDaYun.stem}${displayDaYun.branch}` : undefined,
        detail: `${displayDaYun.start_age}-${displayDaYun.end_age}歲 · ${effectiveFlowYear}流年`,
      });
    }

    if (displayLiuNian) {
      list.push({
        key: 'liu_nian',
        label: '流年',
        palace: displayLiuNian.palace,
        branch: displayLiuNian.branch,
        timeGanzhi: formatTemporalGanzhi(displayLiuNian),
        detail: `${effectiveFlowYear}年`,
      });
    }

    if (displayLiuYue) {
      list.push({
        key: 'liu_yue',
        label: '流月',
        palace: displayLiuYue.palace,
        branch: displayLiuYue.branch,
        timeGanzhi: formatTemporalGanzhi(displayLiuYue),
        detail: `${targetMonth}月`,
      });
    }

    if (displayLiuRi) {
      list.push({
        key: 'liu_ri',
        label: '流日',
        palace: displayLiuRi.palace,
        branch: displayLiuRi.branch,
        timeGanzhi: formatTemporalGanzhi(displayLiuRi),
        detail: `${targetMonth}月${targetDay}日`,
      });
    }

    return list;
  }, [displayDaYun, displayLiuNian, displayLiuYue, displayLiuRi, effectiveFlowYear, targetMonth, targetDay]);

  const keyInsights = useMemo(() => {
    const mingHost = palaceByBranch.get(data.ming_gong);
    const shenHost = palaceByBranch.get(data.shen_gong);
    const originBranch = branchByPalace.get(data.origin_palace);

    const items: Array<{
      label: string;
      value: string;
      palace?: string;
      branch?: string;
      temporalInfo?: {
        stem: string;
        timeBranch?: string;
        ageRange?: string;
        timeLabel?: string;
      };
    }> = [
      // 1. 來因宮
      {
        label: `來因宮（${originBranch || ''}）`,
        value: data.origin_palace,
        palace: data.origin_palace,
        branch: originBranch,
      },
      // 2. 命宮
      {
        label: `命宮（${data.ming_gong}）`,
        value: mingHost || data.ming_gong,
        palace: mingHost,
        branch: data.ming_gong,
      },
      // 3. 身宮
      {
        label: `身宮（${data.shen_gong}）`,
        value: shenHost || data.shen_gong,
        palace: shenHost,
        branch: data.shen_gong,
      },
    ];

    // 4. 大限、5. 流年、6. 流月、7. 流日（可點擊的運限）
    // 使用動態計算的大限數據（如果已加載）
    if (displayDaYun) {
      items.push({
        label: `大限（${displayDaYun.start_age}-${displayDaYun.end_age}歲 · ${effectiveFlowYear}流年）`,
        value: displayDaYun.palace,
        palace: displayDaYun.palace,
        branch: displayDaYun.branch,
        temporalInfo: {
          stem: displayDaYun.stem,
          ageRange: `${displayDaYun.start_age}-${displayDaYun.end_age}歲`,
          timeLabel: `${effectiveFlowYear}流年`,
        },
      });
    }

    if (displayLiuNian) {
      items.push({
        label: `流年（${formatTemporalGanzhi(displayLiuNian)}｜落宮：${displayLiuNian.branch}宮 · ${effectiveFlowYear}年）`,
        value: displayLiuNian.palace,
        palace: displayLiuNian.palace,
        branch: displayLiuNian.branch,
        temporalInfo: {
          stem: displayLiuNian.stem,
          timeBranch: displayLiuNian.time_branch,
          timeLabel: `${effectiveFlowYear}年`,
        },
      });
    }

    if (displayLiuYue) {
      items.push({
        label: `流月（${formatTemporalGanzhi(displayLiuYue)}｜落宮：${displayLiuYue.branch}宮 · ${targetMonth}月）`,
        value: displayLiuYue.palace,
        palace: displayLiuYue.palace,
        branch: displayLiuYue.branch,
        temporalInfo: {
          stem: displayLiuYue.stem,
          timeBranch: displayLiuYue.time_branch,
          timeLabel: `${targetMonth}月`,
        },
      });
    }

    if (displayLiuRi) {
      items.push({
        label: `流日（${formatTemporalGanzhi(displayLiuRi)}｜落宮：${displayLiuRi.branch}宮 · ${targetMonth}月${targetDay}日）`,
        value: displayLiuRi.palace,
        palace: displayLiuRi.palace,
        branch: displayLiuRi.branch,
        temporalInfo: {
          stem: displayLiuRi.stem,
          timeBranch: displayLiuRi.time_branch,
          timeLabel: `${targetMonth}月${targetDay}日`,
        },
      });
    }

    // 8. 納音、9. 年柱、10. 五行局（不可點擊）
    items.push(
      {
        label: '納音',
        value: data.na_yin,
      },
      {
        label: '年柱',
        value: data.year_pillar,
      },
      {
        label: '五行局',
        value: data.wuxing,
      }
    );

    return items;
  }, [data, palaceByBranch, branchByPalace, currentTemporalData, targetMonth, targetDay, effectiveFlowYear]);

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
                <div className="heading-lg">{(userInfo?.name || '未知命主')} 命盤</div>
                <div className="body-sm">{(userInfo?.gender === 'male' ? '乾造' : '坤造')} | {(userInfo?.year || '-')}年{(userInfo?.month || '-')}月{(userInfo?.day || '-')}日{(userInfo?.hour || '-')}時</div>
              </div>
            </div>
            <button className="btn-secondary" onClick={() => setDarkMode(!darkMode)}>
              {darkMode ? <Sun size={16} /> : <Moon size={16} />}
              {darkMode ? '淺色' : '深色'}
            </button>
          </div>

          <div className="metric-grid">
            {keyInsights.map((item) => {
              const isSelected = selectedInsight === item.label;
              return (
                <button
                  key={item.label}
                  className={`metric-card ${isSelected ? 'is-active' : ''}`}
                  style={{
                    cursor: item.palace ? 'pointer' : 'default',
                    textAlign: 'left',
                    border: isSelected ? '2px solid var(--cta)' : 'none',
                    background: isSelected ? 'rgba(178, 135, 70, 0.15)' : 'var(--surface)',
                    boxShadow: isSelected ? '0 0 0 1px var(--cta)' : undefined,
                  }}
                  onClick={() => {
                    if (item.palace) {
                      setSelectedInsight(item.label);
                      const detailLabel = item.temporalInfo
                        ? item.label.split('（')[0]
                        : item.label;
                      setFocusedPalace({
                        label: detailLabel,
                        palaceName: item.palace,
                        branch: item.branch,
                        temporalInfo: item.temporalInfo,
                      });
                    }
                  }}
                  disabled={!item.palace}
                  title={item.palace ? `點擊查看 ${item.label} 的三方四正` : undefined}
                >
                  <div className="metric-label">{item.label}</div>
                  <div className="metric-value">{item.value}</div>
                </button>
              );
            })}
          </div>

          {/* 大限切換控制 */}
          {data.da_yun && data.da_yun.length > 0 && (
            <>
              <div style={{ 
                display: 'flex', 
                alignItems: 'center', 
                gap: '0.75rem', 
                padding: '0.75rem 1rem',
                background: 'rgba(178, 135, 70, 0.08)',
                borderRadius: '0.5rem',
                marginTop: '0.5rem'
              }}>
                <span className="body-sm" style={{ color: 'var(--text-secondary)', whiteSpace: 'nowrap' }}>
                  大限切換：
                </span>
                <button 
                  className="btn-secondary" 
                  onClick={handlePrevDaYun}
                  disabled={currentDaYunIndex === 0 || isTemporalLoading}
                  style={{ padding: '0.375rem 0.75rem', fontSize: '0.875rem' }}
                >
                  ← 上一個
                </button>
                
                <div style={{ 
                  display: 'flex', 
                  gap: '0.375rem', 
                  flexWrap: 'wrap',
                  flex: 1,
                  justifyContent: 'center'
                }}>
                  {data.da_yun.map((dy, index) => (
                    <button
                      key={dy.index}
                      onClick={() => handleDaYunClick(index)}
                      disabled={isTemporalLoading}
                      style={{
                        padding: '0.375rem 0.625rem',
                        fontSize: '0.8rem',
                        borderRadius: '0.375rem',
                        border: currentDaYunIndex === index 
                          ? '2px solid var(--cta)' 
                          : '1px solid var(--border)',
                        background: currentDaYunIndex === index 
                          ? 'rgba(178, 135, 70, 0.2)' 
                          : 'var(--surface)',
                        color: currentDaYunIndex === index 
                          ? 'var(--cta)' 
                          : 'var(--text-secondary)',
                        cursor: isTemporalLoading ? 'not-allowed' : 'pointer',
                        opacity: isTemporalLoading ? 0.6 : 1,
                        minWidth: '3.5rem',
                        textAlign: 'center'
                      }}
                      title={`${dy.start_age}-${dy.end_age}歲 (${dy.palace})`}
                    >
                      {dy.start_age}-{dy.end_age}
                    </button>
                  ))}
                </div>
                
                <button 
                  className="btn-secondary" 
                  onClick={handleNextDaYun}
                  disabled={currentDaYunIndex >= (data.da_yun.length - 1) || isTemporalLoading}
                  style={{ padding: '0.375rem 0.75rem', fontSize: '0.875rem' }}
                >
                  下一個 →
                </button>
                
                {isTemporalLoading && (
                  <span className="body-sm" style={{ color: 'var(--cta)', marginLeft: '0.5rem' }}>
                    計算中...
                  </span>
                )}
              </div>

              <div className="card section-stack" style={{ marginTop: '0.75rem', gap: '0.65rem' }}>
                <div className="heading-sm">流運直觀操作面板</div>

                <div className="body-sm" style={{ color: 'var(--text-secondary)' }}>流年（10 年）</div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                  <button
                    className="btn-secondary"
                    onClick={() => handleFlowYearShift(-1)}
                    disabled={isTemporalLoading}
                    style={{ padding: '0.25rem 0.5rem' }}
                    title="往前一年"
                  >
                    ←
                  </button>
                  <div className="tab-row" style={{ flexWrap: 'wrap', flex: 1 }}>
                    {flowYearOptions.map((year) => {
                      const selected = effectiveFlowYear === year;
                      return (
                        <button
                          key={`flow-year-${year}`}
                          className={`tab-button ${selected ? 'is-active' : ''}`}
                          onClick={() => handleFlowYearSelect(year)}
                          disabled={isTemporalLoading}
                        >
                          {year}
                        </button>
                      );
                    })}
                  </div>
                  <button
                    className="btn-secondary"
                    onClick={() => handleFlowYearShift(1)}
                    disabled={isTemporalLoading}
                    style={{ padding: '0.25rem 0.5rem' }}
                    title="往後一年"
                  >
                    →
                  </button>
                </div>

                <div className="body-sm" style={{ color: 'var(--text-secondary)' }}>流月（12 月）</div>
                <div className="tab-row" style={{ flexWrap: 'wrap', flex: 1 }}>
                  {flowMonthOptions.map((month) => (
                    <button
                      key={`flow-month-${month}`}
                      className={`tab-button ${targetMonth === month ? 'is-active' : ''}`}
                      onClick={() => handleFlowMonthSelect(month)}
                      disabled={isTemporalLoading}
                    >
                      {month}月
                    </button>
                  ))}
                </div>

                <div className="body-sm" style={{ color: 'var(--text-secondary)' }}>流日（12 格，當日居中）</div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                  <button
                    className="btn-secondary"
                    onClick={() => handleFlowDayShift(-1)}
                    disabled={isTemporalLoading}
                    style={{ padding: '0.25rem 0.5rem' }}
                    title="往前12日"
                  >
                    ←
                  </button>
                  <div className="tab-row" style={{ justifyContent: 'center', flexWrap: 'wrap', flex: 1 }}>
                    {flowDayOptions.map((item, idx) => {
                      const selected = item.date.toDateString() === targetDate.toDateString();
                      const emphasis = idx === 5;
                      return (
                        <button
                          key={`flow-day-${idx}-${item.date.toISOString()}`}
                          className={`tab-button ${selected ? 'is-active' : ''}`}
                          onClick={() => handleFlowDaySelect(item.date)}
                          disabled={isTemporalLoading}
                          style={{
                            ...(emphasis ? { borderWidth: '2px' } : {}),
                            opacity: item.isCurrentMonth ? 1 : 0.5,
                          }}
                          title={`${item.month}月${item.day}日${item.isCurrentMonth ? '' : '（跨月）'}`}
                        >
                          {item.month}/{item.day}
                        </button>
                      );
                    })}
                  </div>
                  <button
                    className="btn-secondary"
                    onClick={() => handleFlowDayShift(1)}
                    disabled={isTemporalLoading}
                    style={{ padding: '0.25rem 0.5rem' }}
                    title="往後12日"
                  >
                    →
                  </button>
                </div>
              </div>

            </>
          )}
        </motion.header>

        <section className="section-stack">
          {focusedPalace ? (
            <PalaceDetailView
              palaces={data.palaces}
              focus={focusedPalace}
              temporalLayers={temporalLayers}
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
                    currentDaYun={displayDaYun}
                    daYun={data.da_yun}
                    temporalPositions={temporalPositions}
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
              userInfo={currentUser}
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
