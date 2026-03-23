import React, { useMemo, useState } from 'react';
import { motion } from 'framer-motion';
import { Crown, Sparkles, Star } from 'lucide-react';
import '../styles/design-system.css';

interface TransformData {
  star: string;
  transformation: string;
  display: string;
}

interface FlyStage {
  type: string;
  star: string;
  target: string;
  motive?: string;
  action?: string;
  trap?: string;
}

interface PalaceData {
  branch: string;
  palace_gan?: string;
  stars: string[];
  star_details?: Array<{
    name: string;
    brightness?: string;
  }>;
  assistant_stars?: string[];
  assistant_star_details?: Array<{
    name: string;
    brightness?: string;
  }>;
  secondary_stars?: string[];
  chang_sheng?: string;
  bo_shi?: string;
  natal_transforms?: TransformData[];
  liu_nian_stars?: string[];
  liu_nian_transforms?: TransformData[];
  liu_yue_stars?: string[];
  liu_yue_transforms?: TransformData[];
  liu_ri_stars?: string[];
  liu_ri_transforms?: TransformData[];
  da_yun_ages?: string[];
  fly_hua?: {
    from_palace: string;
    stem: string;
    stages: FlyStage[];
  };
}

interface TemporalPositionSummary {
  key: 'da_yun' | 'liu_nian' | 'liu_yue' | 'liu_ri';
  label: string;
  palace?: string;
  branch?: string;
  timeGanzhi?: string;
  detail?: string;
}

interface DaYunData {
  index: number;
  start_age: number;
  end_age: number;
  stem: string;
  branch: string;
  palace: string;
}

interface ZiweiChartProps {
  palaces: Record<string, PalaceData>;
  mingGong: string;
  shenGong: string;
  originPalace: string;
  currentDaYun?: DaYunData | null;
  daYun?: DaYunData[];
  temporalPositions?: TemporalPositionSummary[];
  className?: string;
}

const branchLayout = [
  { branch: '巳', area: '1 / 1 / 2 / 2' },
  { branch: '午', area: '1 / 2 / 2 / 3' },
  { branch: '未', area: '1 / 3 / 2 / 4' },
  { branch: '申', area: '1 / 4 / 2 / 5' },
  { branch: '酉', area: '2 / 4 / 3 / 5' },
  { branch: '戌', area: '3 / 4 / 4 / 5' },
  { branch: '亥', area: '4 / 4 / 5 / 5' },
  { branch: '子', area: '4 / 3 / 5 / 4' },
  { branch: '丑', area: '4 / 2 / 5 / 3' },
  { branch: '寅', area: '4 / 1 / 5 / 2' },
  { branch: '卯', area: '3 / 1 / 4 / 2' },
  { branch: '辰', area: '2 / 1 / 3 / 2' },
] as const;

const palaceCycle = [
  '命宮', '兄弟宮', '夫妻宮', '子女宮', '財帛宮', '疾厄宮',
  '遷移宮', '僕役宮', '官祿宮', '田宅宮', '福德宮', '父母宮',
];

const resolveShiftedPalace = (sourcePalace: string, newBasePalace: string): string | null => {
  const sourceIndex = palaceCycle.indexOf(sourcePalace);
  const baseIndex = palaceCycle.indexOf(newBasePalace);
  if (sourceIndex === -1 || baseIndex === -1) return null;
  // 從 newBasePalace 的角度看 sourcePalace 是什麼宮
  // 計算相對位置：(base 到 source 的順時針距離)
  return palaceCycle[(baseIndex - sourceIndex + 12) % 12];
};

const mainStars = new Set(['紫微', '天機', '太陽', '武曲', '天同', '廉貞', '天府', '太陰', '貪狼', '巨門', '天相', '天梁', '七殺', '破軍']);

const importantAuxStars = new Set(['左輔', '右弼', '文昌', '文曲', '天魁', '天鉞']);

// 四化表：每個天干對應的化祿、化權、化科、化忌星曜
const transformationTable: Record<string, { lu: string; quan: string; ke: string; ji: string }> = {
  '甲': { lu: '廉貞', quan: '破軍', ke: '武曲', ji: '太陽' },
  '乙': { lu: '天機', quan: '天梁', ke: '紫微', ji: '太陰' },
  '丙': { lu: '天同', quan: '天機', ke: '文昌', ji: '廉貞' },
  '丁': { lu: '太陰', quan: '天同', ke: '天機', ji: '巨門' },
  '戊': { lu: '貪狼', quan: '太陰', ke: '右弼', ji: '天機' },
  '己': { lu: '武曲', quan: '貪狼', ke: '天梁', ji: '文曲' },
  '庚': { lu: '太陽', quan: '武曲', ke: '太陰', ji: '天同' },
  '辛': { lu: '巨門', quan: '太陽', ke: '文曲', ji: '文昌' },
  '壬': { lu: '天梁', quan: '紫微', ke: '左輔', ji: '武曲' },
  '癸': { lu: '破軍', quan: '巨門', ke: '太陰', ji: '貪狼' },
};

// 計算宮位飛星四化：根據宮干返回該宮位星曜的化祿、化權、化科、化忌
const calculatePalaceFlyHua = (palaceGan: string, starsInPalace: string[]): Record<string, string> => {
  const result: Record<string, string> = {};
  const transforms = transformationTable[palaceGan];
  if (!transforms) return result;

  starsInPalace.forEach((star) => {
    if (star === transforms.lu) result[star] = '化祿';
    else if (star === transforms.quan) result[star] = '化權';
    else if (star === transforms.ke) result[star] = '化科';
    else if (star === transforms.ji) result[star] = '化忌';
  });

  return result;
};

const ZiweiChart: React.FC<ZiweiChartProps> = ({
  palaces,
  mingGong,
  shenGong,
  originPalace,
  currentDaYun,
  daYun: _daYun = [],
  temporalPositions = [],
  className = ''
}) => {
  const palaceEntries = useMemo(() => {
    return Object.entries(palaces).map(([palaceName, data]) => ({
      palaceName,
      ...data,
    }));
  }, [palaces]);

  const branchToEntry = useMemo(() => {
    const map = new Map<string, PalaceData & { palaceName: string }>();
    palaceEntries.forEach((entry) => {
      map.set(entry.branch, entry);
    });
    return map;
  }, [palaceEntries]);

  const [selectedPalaceName, setSelectedPalaceName] = useState<string>(originPalace || palaceEntries[0]?.palaceName || '命宮');

  const selectedPalaceGan = useMemo(() => {
    const entry = palaceEntries.find(e => e.palaceName === selectedPalaceName);
    return entry?.palace_gan || '';
  }, [selectedPalaceName, palaceEntries]);

  const mappingBasePalace = selectedPalaceName;

  const selectedPalaceFlyHua = useMemo(() => {
    if (!selectedPalaceGan) return {} as Record<string, string>;
    const allStarsInSystem = new Set<string>();
    palaceEntries.forEach(entry => {
      (entry.stars || []).forEach(star => allStarsInSystem.add(star));
      (entry.assistant_stars || []).forEach(star => allStarsInSystem.add(star));
    });
    return calculatePalaceFlyHua(selectedPalaceGan, Array.from(allStarsInSystem));
  }, [selectedPalaceGan, palaceEntries]);

  const formatTemporalBadge = (item: TemporalPositionSummary): string => {
    const timeText = item.timeGanzhi ? `${item.label} ${item.timeGanzhi}` : item.label;
    const landingText = item.branch ? `｜落宮：${item.branch}宮` : '';
    return `${timeText}${landingText}`;
  };

  const selectedEntry = palaceEntries.find((entry) => entry.palaceName === selectedPalaceName) ?? palaceEntries[0];
  const selectedIndex = palaceCycle.indexOf(selectedEntry?.palaceName || '命宮');
  const sanFangSet = new Set<string>(
    selectedIndex >= 0
      ? [
          palaceCycle[selectedIndex],
          palaceCycle[(selectedIndex + 4) % 12],
          palaceCycle[(selectedIndex + 6) % 12],
          palaceCycle[(selectedIndex + 8) % 12],
        ]
      : [],
  );

  const temporalHighlights = new Map<string, string[]>();
  const pushHighlight = (palaceName?: string, label?: string) => {
    if (!palaceName || !label) return;
    temporalHighlights.set(palaceName, [...(temporalHighlights.get(palaceName) || []), label]);
  };
  temporalPositions.forEach((item) => {
    pushHighlight(item.palace, item.label);
  });


  const selectedStarList: Array<{ name: string; brightness?: string }> =
    selectedEntry?.star_details && selectedEntry.star_details.length > 0
      ? selectedEntry.star_details
      : (selectedEntry?.stars || []).map((star) => ({ name: star }));

  return (
    <div className={`ziwei-chart-container ${className}`}>
      <motion.div className="section-stack" initial={{ opacity: 0, y: -16 }} animate={{ opacity: 1, y: 0 }}>
        <div className="card" style={{ width: '100%' }}>
          <div
            className="palace-grid"
            style={{
              display: 'grid',
              gridTemplateColumns: 'repeat(4, minmax(0, 1fr))',
              gridTemplateRows: 'repeat(4, minmax(120px, auto))',
              width: '100%',
            }}
          >
              {branchLayout.map(({ branch, area }) => {
                const entry = branchToEntry.get(branch);
                if (!entry) return null;

                const highlightLabels = temporalHighlights.get(entry.palaceName) || [];
                const isSelected = entry.palaceName === selectedPalaceName;
                const isSanFang = sanFangSet.has(entry.palaceName);
                const shiftedPalaceName = resolveShiftedPalace(mappingBasePalace, entry.palaceName);

                // 三方四正邊框顏色邏輯
                const isOpposite = entry.palaceName === palaceCycle[(selectedIndex + 6) % 12];
                const isSanHe1 = entry.palaceName === palaceCycle[(selectedIndex + 4) % 12];
                const isSanHe2 = entry.palaceName === palaceCycle[(selectedIndex + 8) % 12];

                let borderStyle: React.CSSProperties = {};
                if (isSelected) {
                  borderStyle = { border: '2px solid var(--cta)' };
                } else if (isOpposite) {
                  borderStyle = { border: '2px solid #dc2626' };
                } else if (isSanHe1 || isSanHe2) {
                  borderStyle = { border: '2px solid #2563eb' };
                } else if (isSanFang) {
                  borderStyle = { border: '2px solid #8b5cf6' };
                }

                const classNames = ['palace-cell'];
                if (entry.palaceName === originPalace) classNames.push('palace-cell-origin');
                if (entry.branch === mingGong) classNames.push('palace-cell-current');
                if (isSelected) classNames.push('is-selected');
                if (isSanFang) classNames.push('is-related');

                return (
                  <button
                    key={branch}
                    className={classNames.join(' ')}
                    style={{ 
                      gridArea: area, 
                      textAlign: 'left', 
                      display: 'flex',
                      flexDirection: 'column',
                      ...borderStyle 
                    }}
                    onClick={() => setSelectedPalaceName(entry.palaceName)}
                  >
                    <div className="palace-cell-header" style={{ display: 'flex', justifyContent: 'space-between', gap: '0.5rem', alignItems: 'flex-start' }}>
                      <div style={{ display: 'flex', alignItems: 'baseline', gap: '0.35rem' }}>
                        <span style={{ fontWeight: 700, fontSize: '0.95rem' }}>{entry.palaceName}</span>
                        <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)', fontWeight: 400 }}>→{shiftedPalaceName || '—'}</span>
                      </div>
                      <div style={{ display: 'flex', gap: '0.35rem', flexWrap: 'wrap', justifyContent: 'flex-end' }}>
                        {entry.palaceName === originPalace && <Crown size={14} color="var(--cta)" />}
                        {entry.branch === shenGong && <span className="tab-button" style={{ fontSize: '0.6rem', padding: '0.1rem 0.3rem' }}>身</span>}
                      </div>
                    </div>

                    <div className="palace-cell-stars" style={{ display: 'grid', gap: '0.2rem', fontSize: '0.9rem', flex: 1 }}>
                      {(() => {
                        const palaceGan = entry.palace_gan || '';
                        const allStars = entry.star_details && entry.star_details.length > 0
                          ? entry.star_details.map(s => s.name)
                          : entry.stars;
                        const flyHua = palaceGan ? calculatePalaceFlyHua(palaceGan, allStars) : {};
                        const starToTransform = new Map<string, string>();
                        (entry.natal_transforms || []).forEach(t => {
                          starToTransform.set(t.star, t.transformation);
                        });

                        const starList = entry.star_details && entry.star_details.length > 0
                          ? entry.star_details
                          : entry.stars.map((star) => ({ name: star } as { name: string; brightness?: string }));

                        // 輔星和煞星現在從 assistant_star_details 獲取，包含亮度資訊
                        const assistantStarList = entry.assistant_star_details && entry.assistant_star_details.length > 0
                          ? entry.assistant_star_details
                          : (entry.assistant_stars || []).map((star) => ({ name: star } as { name: string; brightness?: string }));

                        const mainStarList = starList.filter(s => mainStars.has(s.name));
                        const importantAuxFromStars = starList.filter(s => importantAuxStars.has(s.name));
                        const importantAuxFromAssistant = assistantStarList
                          .filter(s => importantAuxStars.has(s.name));
                        const importantStars = [...mainStarList, ...importantAuxFromStars, ...importantAuxFromAssistant];
                        const otherAuxStars = assistantStarList.filter(s => !importantAuxStars.has(s.name));

                        return (
                          <>
                            {importantStars.map((star) => {
                              const localFlyTransformation = flyHua[star.name];
                              const natalTransformation = starToTransform.get(star.name);
                              const selectedFlyTransformation = selectedPalaceFlyHua[star.name];
                              const isMainStar = mainStars.has(star.name);
                              const isSelectedPalace = entry.palaceName === selectedPalaceName;
                              return (
                                <div key={`${entry.palaceName}-${star.name}`} style={{ display: 'flex', alignItems: 'center', gap: '0.3rem' }}>
                                  {isMainStar ? <Star size={12} color="var(--cta)" /> : <Sparkles size={12} color="var(--secondary)" />}
                                  <span style={{
                                    fontWeight: isMainStar ? 600 : 500,
                                    color: isMainStar ? 'var(--primary)' : 'var(--secondary)'
                                  }}>
                                    {star.name}
                                  </span>
                                  {star.brightness && (
                                    <span style={{ color: 'var(--text-muted)', fontSize: '0.75rem' }}>
                                      {star.brightness}
                                    </span>
                                  )}
                                  {natalTransformation && (
                                    <span style={{
                                      color: 'var(--accent)',
                                      fontSize: '0.75rem',
                                      fontWeight: 600,
                                      background: 'rgba(163, 63, 47, 0.1)',
                                      padding: '0 0.2rem',
                                      borderRadius: '0.15rem',
                                      whiteSpace: 'nowrap'
                                    }}>
                                      {natalTransformation}
                                    </span>
                                  )}
                                  {isSelectedPalace && localFlyTransformation && (
                                    <span style={{
                                      color: 'var(--cta)',
                                      fontSize: '0.75rem',
                                      fontWeight: 500,
                                      whiteSpace: 'nowrap'
                                    }}>
                                      飛{localFlyTransformation}
                                    </span>
                                  )}
                                  {!isSelectedPalace && selectedFlyTransformation && (
                                    <span style={{
                                      color: '#7c3aed',
                                      fontSize: '0.75rem',
                                      fontWeight: 600,
                                      background: 'rgba(124, 58, 237, 0.1)',
                                      padding: '0 0.2rem',
                                      borderRadius: '0.15rem',
                                      whiteSpace: 'nowrap'
                                    }}>
                                      {selectedPalaceName}{selectedFlyTransformation}
                                    </span>
                                  )}
                                </div>
                              );
                            })}
                            {otherAuxStars.length > 0 && (
                              <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.2rem', marginTop: '0.15rem' }}>
                                {otherAuxStars.map((star) => (
                                  <span key={`${entry.palaceName}-assistant-${star.name}`} style={{ color: 'var(--secondary)', fontSize: '0.8rem' }}>
                                    {star.name}{star.brightness ? ` ${star.brightness}` : ''}
                                  </span>
                                ))}
                              </div>
                            )}
                            {entry.secondary_stars && entry.secondary_stars.length > 0 && (
                              <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.2rem' }}>
                                {entry.secondary_stars.map((star) => (
                                  <span key={`${entry.palaceName}-secondary-${star}`} style={{ color: 'var(--text-soft)', fontSize: '0.75rem' }}>{star}</span>
                                ))}
                              </div>
                            )}
                          </>
                        );
                      })()}
                    </div>

{highlightLabels.length > 0 && (
                      <div style={{ marginTop: '0.5rem', display: 'flex', flexWrap: 'wrap', gap: '0.25rem', justifyContent: 'flex-end' }}>
                        {highlightLabels.map((label) => (
                          <span key={`${entry.palaceName}-${label}`} className="tab-button" style={{ fontSize: '0.7rem', padding: '0.15rem 0.4rem' }}>{label}</span>
                        ))}
                      </div>
                    )}

                    <div style={{ marginTop: 'auto', paddingTop: '0.5rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', fontSize: '0.72rem', color: 'var(--text-soft)', borderTop: '1px solid rgba(178, 135, 70, 0.1)', width: '100%' }}>
                      <span style={{ flex: 1, textAlign: 'left' }}>{entry.palace_gan || '—'}·{entry.branch}宮</span>
                      {entry.da_yun_ages && entry.da_yun_ages.length > 0 && (
                        <span style={{ flex: 1, textAlign: 'center', color: 'var(--secondary)', fontWeight: 500 }}>{entry.da_yun_ages.join('、')}</span>
                      )}
                      <span style={{ flex: 1, textAlign: 'right' }}>{entry.chang_sheng || ''}</span>
                    </div>
                  </button>
                );
              })}

              <div
                className="card card-gold"
                style={{
                  gridArea: '2 / 2 / 4 / 4',
                  borderRadius: '1.2rem',
                  display: 'flex',
                  flexDirection: 'column',
                  justifyContent: 'center',
                  gap: '0.75rem',
                  textAlign: 'center',
                  padding: '1rem',
                }}
              >
                <div className="heading-md" style={{ fontSize: '1.6rem' }}>{selectedEntry?.palaceName}</div>
                <div className="body-sm">{selectedEntry?.branch}位 · 宮干 {selectedEntry?.palace_gan || '—'}</div>
                <div className="tab-row" style={{ justifyContent: 'center' }}>
                  {selectedStarList.slice(0, 4).map((item) => (
                    <span key={`center-${item.name}`} className="tab-button">
                      {item.name}{item.brightness ? ` ${item.brightness}` : ''}
                    </span>
                  ))}
                </div>
                <div className="tab-row" style={{ justifyContent: 'center' }}>
                  {currentDaYun && <span className="tab-button">大限 {currentDaYun.start_age}-{currentDaYun.end_age}</span>}
                  {temporalPositions.filter((item) => item.key !== 'da_yun').map((item) => (
                    <button
                      key={`center-temporal-${item.key}`}
                      className="tab-button"
                      onClick={() => item.palace && setSelectedPalaceName(item.palace)}
                      title={item.detail || item.label}
                    >
                      {formatTemporalBadge(item)}
                    </button>
                  ))}
                </div>
                <div className="body-sm">來因宮：{originPalace} · 命宮地支：{mingGong} · 身宮地支：{shenGong}</div>
                {selectedPalaceGan && (
                  <div style={{ 
                    marginTop: '0.5rem', 
                    padding: '0.5rem', 
                    background: 'rgba(124, 58, 237, 0.08)', 
                    borderRadius: '0.5rem',
                    border: '1px solid rgba(124, 58, 237, 0.2)'
                  }}>
                    <div style={{ fontSize: '0.75rem', color: '#7c3aed', marginBottom: '0.25rem' }}>
                      {selectedPalaceName}宮干{selectedPalaceGan}飛化四化
                    </div>
                    <div className="tab-row" style={{ justifyContent: 'center' }}>
                      {Object.entries(selectedPalaceFlyHua).slice(0, 4).map(([star, trans]) => (
                        <span key={`fly-${star}`} className="tab-button" style={{ 
                          fontSize: '0.7rem', 
                          padding: '0.15rem 0.4rem',
                          background: 'rgba(124, 58, 237, 0.15)',
                          color: '#7c3aed',
                          border: '1px solid rgba(124, 58, 237, 0.3)'
                        }}>
                          {star}{trans}
                        </span>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
      </motion.div>
    </div>
  );
};

export default ZiweiChart;
