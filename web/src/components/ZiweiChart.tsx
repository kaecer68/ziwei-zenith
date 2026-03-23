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
  return palaceCycle[(sourceIndex - baseIndex + 12) % 12];
};

const mainStars = new Set(['紫微', '天機', '太陽', '武曲', '天同', '廉貞', '天府', '太陰', '貪狼', '巨門', '天相', '天梁', '七殺', '破軍']);

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

  const formatTemporalBadge = (item: TemporalPositionSummary): string => {
    const timeText = item.timeGanzhi ? `${item.label} ${item.timeGanzhi}` : item.label;
    const landingText = item.branch ? `｜落宮：${item.branch}宮` : '';
    return `${timeText}${landingText}`;
  };

  const selectedEntry = palaceEntries.find((entry) => entry.palaceName === selectedPalaceName) ?? palaceEntries[0];
  const mappingBasePalace = '命宮';
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
        <div className="insight-grid" style={{ alignItems: 'start' }}>
          <div className="card">
            <div
              className="palace-grid"
              style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(4, minmax(0, 1fr))',
                gridTemplateRows: 'repeat(4, minmax(120px, auto))',
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
                    style={{ gridArea: area, textAlign: 'left', ...borderStyle }}
                    onClick={() => setSelectedPalaceName(entry.palaceName)}
                  >
                    <div className="palace-cell-header" style={{ display: 'flex', justifyContent: 'space-between', gap: '0.5rem', alignItems: 'flex-start' }}>
                      <div>
                        <div style={{ fontSize: '0.72rem', color: 'var(--text-soft)' }}>{entry.branch}宮 · {entry.palace_gan || ''}</div>
                        <div>{entry.palaceName}</div>
                        <div style={{ fontSize: '0.7rem', color: 'var(--text-muted)' }}>
                           →{shiftedPalaceName || '—'}
                         </div>
                      </div>
                      <div style={{ display: 'flex', gap: '0.35rem', flexWrap: 'wrap', justifyContent: 'flex-end' }}>
                        {entry.palaceName === originPalace && <Crown size={14} color="var(--cta)" />}
                        {entry.branch === mingGong && <span className="tab-button is-active">命</span>}
                        {entry.branch === shenGong && <span className="tab-button">身</span>}
                      </div>
                    </div>

                    <div className="palace-cell-stars" style={{ display: 'grid', gap: '0.2rem' }}>
                      {(() => {
                        const palaceGan = entry.palace_gan || '';
                        const allStars = entry.star_details && entry.star_details.length > 0
                          ? entry.star_details.map(s => s.name)
                          : entry.stars;
                        const flyHua = palaceGan ? calculatePalaceFlyHua(palaceGan, allStars) : {};

                        return (entry.star_details && entry.star_details.length > 0
                          ? entry.star_details
                          : entry.stars.map((star) => ({ name: star } as { name: string; brightness?: string }))
                        ).map((star) => {
                          const transformation = flyHua[star.name];
                          return (
                            <div key={`${entry.palaceName}-${star.name}`} style={{ display: 'flex', alignItems: 'center', gap: '0.35rem', justifyContent: 'space-between' }}>
                              <span style={{ display: 'flex', alignItems: 'center', gap: '0.35rem' }}>
                                {mainStars.has(star.name) ? <Star size={12} color="var(--cta)" /> : <Sparkles size={12} color="var(--text-soft)" />}
                                <span>
                                  {star.name}
                                  {transformation && (
                                    <span style={{ color: 'var(--accent)', fontSize: '0.7rem', marginLeft: '0.2rem' }}>
                                      【{entry.palaceName}{transformation}】
                                    </span>
                                  )}
                                </span>
                              </span>
                              {star.brightness ? <span style={{ color: 'var(--accent)', fontSize: '0.72rem' }}>{star.brightness}</span> : null}
                            </div>
                          );
                        });
                      })()}
                      {(entry.assistant_stars || []).map((star) => (
                        <div key={`${entry.palaceName}-assistant-${star}`} style={{ color: 'var(--secondary)' }}>{star}</div>
                      ))}
                      {(entry.secondary_stars || []).map((star) => (
                        <div key={`${entry.palaceName}-secondary-${star}`} style={{ color: 'var(--text-soft)' }}>{star}</div>
                      ))}
                    </div>

                    {(entry.natal_transforms?.length || highlightLabels.length || entry.da_yun_ages?.length) ? (
                      <div style={{ marginTop: '0.65rem', display: 'grid', gap: '0.35rem' }}>
                        <div className="tab-row">
                          {(entry.natal_transforms || []).map((item) => (
                            <span key={item.display} className="tab-button is-active">{item.display}</span>
                          ))}
                          {highlightLabels.map((label) => (
                            <span key={`${entry.palaceName}-${label}`} className="tab-button">{label}</span>
                          ))}
                        </div>
                        {entry.da_yun_ages && entry.da_yun_ages.length > 0 && (
                          <div className="body-sm" style={{ color: 'var(--secondary)' }}>大限：{entry.da_yun_ages.join('、')}</div>
                        )}
                      </div>
                    ) : null}

                    {(entry.chang_sheng || entry.bo_shi) && (
                      <div className="body-sm" style={{ marginTop: '0.45rem', color: 'var(--text-muted)' }}>
                        {entry.chang_sheng ? `${entry.chang_sheng}` : ''}{entry.chang_sheng && entry.bo_shi ? ' · ' : ''}{entry.bo_shi ? `${entry.bo_shi}` : ''}
                      </div>
                    )}
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
              </div>
            </div>
          </div>
        </div>
      </motion.div>
    </div>
  );
};

export default ZiweiChart;
