import React, { useMemo, useState } from 'react';
import { motion } from 'framer-motion';
import { Crown, Orbit, Sparkles, Star } from 'lucide-react';
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

type MappingPerspective = 'natal' | 'focused';

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

const ZiweiChart: React.FC<ZiweiChartProps> = ({
  palaces,
  mingGong,
  shenGong,
  originPalace,
  currentDaYun,
  daYun = [],
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
  const [mappingPerspective, setMappingPerspective] = useState<MappingPerspective>('natal');

  const formatTemporalBadge = (item: TemporalPositionSummary): string => {
    const timeText = item.timeGanzhi ? `${item.label} ${item.timeGanzhi}` : item.label;
    const landingText = item.branch ? `｜落宮：${item.branch}宮` : '';
    return `${timeText}${landingText}`;
  };

  const selectedEntry = palaceEntries.find((entry) => entry.palaceName === selectedPalaceName) ?? palaceEntries[0];
  const mappingBasePalace = mappingPerspective === 'natal' ? '命宮' : (selectedEntry?.palaceName || '命宮');
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

  const relationshipSummary = [
    { label: '本宮', palace: palaceCycle[selectedIndex] },
    { label: '三合', palace: palaceCycle[(selectedIndex + 4) % 12] },
    { label: '對宮', palace: palaceCycle[(selectedIndex + 6) % 12] },
    { label: '三合', palace: palaceCycle[(selectedIndex + 8) % 12] },
  ]
    .filter((item) => item.palace)
    .map((item) => ({
      ...item,
      mapped: item.palace ? resolveShiftedPalace(mappingBasePalace, item.palace) : null,
    }));

  const transformSections = [
    { label: '本命四化', items: selectedEntry?.natal_transforms || [] },
    { label: '流年四化', items: selectedEntry?.liu_nian_transforms || [] },
    { label: '流月四化', items: selectedEntry?.liu_yue_transforms || [] },
    { label: '流日四化', items: selectedEntry?.liu_ri_transforms || [] },
  ];

  const selectedStarList: Array<{ name: string; brightness?: string }> =
    selectedEntry?.star_details && selectedEntry.star_details.length > 0
      ? selectedEntry.star_details
      : (selectedEntry?.stars || []).map((star) => ({ name: star }));

  return (
    <div className={`ziwei-chart-container ${className}`}>
      <motion.div className="section-stack" initial={{ opacity: 0, y: -16 }} animate={{ opacity: 1, y: 0 }}>
        <div>
          <h2 className="heading-lg">紫微命盤</h2>
          <p className="body-sm">以地支固定方位排盤，宮名依命宮流轉。點擊宮位可檢視三方四正、四化與流運疊盤。</p>
        </div>

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
                const classNames = ['palace-cell'];
                if (entry.palaceName === originPalace) classNames.push('palace-cell-origin');
                if (entry.branch === mingGong) classNames.push('palace-cell-current');
                if (isSelected) classNames.push('is-selected');
                if (isSanFang) classNames.push('is-related');

                return (
                  <button
                    key={branch}
                    className={classNames.join(' ')}
                    style={{ gridArea: area, textAlign: 'left' }}
                    onClick={() => setSelectedPalaceName(entry.palaceName)}
                  >
                    <div className="palace-cell-header" style={{ display: 'flex', justifyContent: 'space-between', gap: '0.5rem', alignItems: 'flex-start' }}>
                      <div>
                        <div style={{ fontSize: '0.72rem', color: 'var(--text-soft)' }}>{entry.branch}宮 · {entry.palace_gan || ''}</div>
                        <div>{entry.palaceName}</div>
                        <div style={{ fontSize: '0.7rem', color: 'var(--text-muted)' }}>
                          改變對應：{mappingBasePalace}→{shiftedPalaceName || '—'}
                        </div>
                      </div>
                      <div style={{ display: 'flex', gap: '0.35rem', flexWrap: 'wrap', justifyContent: 'flex-end' }}>
                        {entry.palaceName === originPalace && <Crown size={14} color="var(--cta)" />}
                        {entry.branch === mingGong && <span className="tab-button is-active">命</span>}
                        {entry.branch === shenGong && <span className="tab-button">身</span>}
                      </div>
                    </div>

                    <div className="palace-cell-stars" style={{ display: 'grid', gap: '0.2rem' }}>
                      {(entry.star_details && entry.star_details.length > 0
                        ? entry.star_details
                        : entry.stars.map((star) => ({ name: star } as { name: string; brightness?: string }))
                      ).map((star) => (
                        <div key={`${entry.palaceName}-${star.name}`} style={{ display: 'flex', alignItems: 'center', gap: '0.35rem', justifyContent: 'space-between' }}>
                          <span style={{ display: 'flex', alignItems: 'center', gap: '0.35rem' }}>
                            {mainStars.has(star.name) ? <Star size={12} color="var(--cta)" /> : <Sparkles size={12} color="var(--text-soft)" />}
                            <span>{star.name}</span>
                          </span>
                          {star.brightness ? <span style={{ color: 'var(--accent)', fontSize: '0.72rem' }}>{star.brightness}</span> : null}
                        </div>
                      ))}
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
                <div className="body-sm" style={{ color: 'var(--text-soft)' }}>
                  目前視角：{mappingPerspective === 'natal' ? '命宮視角' : '點選宮視角'}（基準宮：{mappingBasePalace}）
                </div>
              </div>
            </div>
          </div>

          <div className="section-stack">
            <div className="card section-stack">
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <Orbit size={18} color="var(--accent)" />
                <div className="heading-md">點宮位後的聯動資訊</div>
              </div>
              <div className="tab-row">
                <button
                  className={`tab-button ${mappingPerspective === 'natal' ? 'is-active' : ''}`}
                  onClick={() => setMappingPerspective('natal')}
                >
                  命宮視角
                </button>
                <button
                  className={`tab-button ${mappingPerspective === 'focused' ? 'is-active' : ''}`}
                  onClick={() => setMappingPerspective('focused')}
                >
                  點選宮視角
                </button>
              </div>
              <div className="tab-row">
                {relationshipSummary.map((item, index) => (
                  <span key={`${item.label}-${index}`} className={`tab-button ${item.palace === selectedEntry?.palaceName ? 'is-active' : ''}`}>
                    {item.label}：{item.palace}（改變對應：{mappingBasePalace}→{item.mapped || '—'}）
                  </span>
                ))}
              </div>
              <div className="body-sm">
                選中 `{selectedEntry?.palaceName}` 後，盤面會同步標示本宮、對宮與兩個三合宮，並疊加大限、流年、流月、流日所在宮位。
              </div>
            </div>

            <div className="card section-stack">
              <div className="heading-md">飛星四化</div>
              {selectedEntry?.fly_hua?.stages?.length ? (
                <div className="section-stack" style={{ gap: '0.75rem' }}>
                  {selectedEntry.fly_hua.stages.map((stage) => (
                    <div key={`${stage.type}-${stage.star}`} className="metric-card">
                      <div className="metric-label">{stage.type}</div>
                      <div className="body-md">{stage.star} → {stage.target || '未定位'}</div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="empty-state">此宮目前沒有可顯示的飛星四化資料。</div>
              )}
            </div>

            <div className="card section-stack">
              <div className="heading-md">四化與流運疊盤</div>
              {transformSections.some((section) => section.items.length > 0) ? (
                transformSections.map((section) => (
                  <div key={section.label}>
                    <div className="heading-sm" style={{ marginBottom: '0.5rem' }}>{section.label}</div>
                    <div className="tab-row">
                      {section.items.length > 0 ? (
                        section.items.map((item) => <span key={`${section.label}-${item.display}`} className="tab-button">{item.display}</span>)
                      ) : (
                        <span className="body-sm">無</span>
                      )}
                    </div>
                  </div>
                ))
              ) : (
                <div className="empty-state">此宮沒有疊入本命 / 流運四化。</div>
              )}
            </div>

            <div className="card section-stack">
              <div className="heading-md">流運定位</div>
              <div className="metric-grid" style={{ gridTemplateColumns: 'repeat(2, minmax(0, 1fr))' }}>
                {temporalPositions.map((item) => (
                  <button
                    key={`position-${item.key}`}
                    className="metric-card"
                    style={{ textAlign: 'left', cursor: item.palace ? 'pointer' : 'default' }}
                    onClick={() => item.palace && setSelectedPalaceName(item.palace)}
                    disabled={!item.palace}
                    title={item.palace ? `點擊定位到${item.palace}` : undefined}
                  >
                    <div className="metric-label">{item.label}</div>
                    <div className="body-md">{item.palace || '未定位'}</div>
                    <div className="body-sm" style={{ color: 'var(--text-soft)' }}>
                      {item.timeGanzhi || '未知'}｜落宮：{item.branch ? `${item.branch}宮` : '未知'}
                    </div>
                    {item.detail && (
                      <div className="body-sm" style={{ color: 'var(--text-muted)' }}>{item.detail}</div>
                    )}
                  </button>
                ))}
              </div>
              {daYun.length > 0 && (
                <div className="body-sm" style={{ color: 'var(--text-muted)' }}>
                  全部大限：{daYun.map((item) => `${item.start_age}-${item.end_age} ${item.palace}`).join(' ｜ ')}
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
