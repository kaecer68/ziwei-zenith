import React, { useMemo } from 'react';
import { ChevronLeft } from 'lucide-react';

interface StarDetail {
  name: string;
  brightness?: string;
  isMain?: boolean;
}

interface TransformInfo {
  star: string;
  transformation: string;
  display: string;
  source?: string;
}

interface PalaceData {
  branch: string;
  palace_gan?: string;
  stars: string[];
  star_details?: StarDetail[];
  assistant_stars?: string[];
  secondary_stars?: string[];
  natal_transforms?: TransformInfo[];
  dayun_transforms?: TransformInfo[];
  liunian_transforms?: TransformInfo[];
  liuyue_transforms?: TransformInfo[];
  liuri_transforms?: TransformInfo[];
}

interface TemporalInfo {
  stem: string;
  timeBranch?: string;
  ageRange?: string;
  timeLabel?: string;
}

interface FocusedPalace {
  label: string;
  palaceName: string;
  branch?: string;
  temporalInfo?: TemporalInfo;
}

interface TemporalLayerInfo {
  label: string;
  palace?: string;
  branch?: string;
  stem?: string;
  timeBranch?: string;
  timeLabel?: string;
}

interface PalaceDetailViewProps {
  palaces: Record<string, PalaceData>;
  focus: FocusedPalace;
  temporalLayers?: TemporalLayerInfo[];
  onBack: () => void;
}

const transformationTable: Record<string, [string, string, string, string]> = {
  '甲': ['廉貞', '破軍', '武曲', '太陽'],
  '乙': ['天機', '天梁', '紫微', '太陰'],
  '丙': ['天同', '天機', '文昌', '廉貞'],
  '丁': ['太陰', '天同', '天機', '巨門'],
  '戊': ['貪狼', '太陰', '右弼', '天機'],
  '己': ['武曲', '貪狼', '天梁', '文曲'],
  '庚': ['太陽', '武曲', '天府', '天同'],
  '辛': ['巨門', '太陽', '文曲', '文昌'],
  '壬': ['天梁', '紫微', '左輔', '武曲'],
  '癸': ['破軍', '巨門', '太陰', '貪狼'],
};

const transformationTypes = ['祿', '權', '科', '忌'];

const mainStarSet = new Set([
  '紫微', '天機', '太陽', '武曲', '天同', '廉貞',
  '天府', '太陰', '貪狼', '巨門', '天相', '天梁', '七殺', '破軍'
]);

interface FlyHuaResult {
  type: string;
  star: string;
  targetPalace: string;
  targetBranch: string;
}

const palaceOrder = [
  '命宮', '父母宮', '福德宮', '田宅宮', '官祿宮', '交友宮',
  '遷移宮', '疾厄宮', '財帛宮', '子女宮', '夫妻宮', '兄弟宮'
];

function normalizePalaceName(palaceName: string): string {
  if (palaceName === '僕役宮') return '交友宮';
  return palaceName;
}

function resolvePalaceEntry(
  palaces: Record<string, PalaceData>,
  palaceName: string
): { key: string; entry: PalaceData } | null {
  const normalized = normalizePalaceName(palaceName);
  const candidates = Array.from(new Set([
    palaceName,
    normalized,
    normalized === '交友宮' ? '僕役宮' : normalized,
  ]));

  for (const key of candidates) {
    const entry = palaces[key];
    if (entry) return { key, entry };
  }

  return null;
}

function resolveShiftedPalace(sourcePalace: string, newBasePalace: string): string | null {
  const source = normalizePalaceName(sourcePalace);
  const base = normalizePalaceName(newBasePalace);
  const sourceIndex = palaceOrder.indexOf(source);
  const baseIndex = palaceOrder.indexOf(base);
  if (sourceIndex === -1 || baseIndex === -1) return null;
  return palaceOrder[(sourceIndex - baseIndex + 12) % 12];
}

function calculateFlyHua(
  palaceGan: string,
  palaces: Record<string, PalaceData>
): FlyHuaResult[] {
  const table = transformationTable[palaceGan];
  if (!table) return [];

  const results: FlyHuaResult[] = [];

  for (let i = 0; i < 4; i++) {
    const starName = table[i];
    const transType = transformationTypes[i];

    for (const [palaceName, data] of Object.entries(palaces)) {
      const allStars = [
        ...(data.stars || []),
        ...(data.assistant_stars || []),
        ...(data.secondary_stars || []),
      ];

      if (allStars.includes(starName)) {
        results.push({
          type: transType,
          star: starName,
          targetPalace: palaceName,
          targetBranch: data.branch,
        });
        break;
      }
    }
  }

  return results;
}

function calculateSanFang(palaceName: string): string[] {
  const normalizedPalace = normalizePalaceName(palaceName);
  const idx = palaceOrder.indexOf(normalizedPalace);
  if (idx === -1) return [normalizedPalace];

  const sanHe1 = (idx + 4) % 12;
  const sanHe2 = (idx + 8) % 12;
  const duiGong = (idx + 6) % 12;

  return [
    normalizedPalace,
    palaceOrder[sanHe1],
    palaceOrder[sanHe2],
    palaceOrder[duiGong],
  ];
}

function getRoleLabel(palaceName: string, focusPalaceName: string, sanFangPalaces: string[]): string {
  const normalizedPalace = normalizePalaceName(palaceName);
  const normalizedFocus = normalizePalaceName(focusPalaceName);
  if (normalizedPalace === normalizedFocus) return '本宮';
  const idx = sanFangPalaces.indexOf(normalizedPalace);
  if (idx === 1 || idx === 2) return '三合宮';
  if (idx === 3) return '對宮';
  return '';
}

function getTransformSourceLabel(source?: string, focusLabel?: string): string {
  switch (source) {
    case 'natal': return '生年';
    case 'dayun': return '大限';
    case 'liunian':
    case 'liunian-fly': return '流年';
    case 'liuyue':
    case 'liuyue-fly': return '流月';
    case 'liuri':
    case 'liuri-fly': return '流日';
    case 'natal-fly':
      // 本命飛星使用宮位名稱
      if (focusLabel?.includes('命宮')) return '命宮';
      if (focusLabel?.includes('身宮')) return '身宮';
      if (focusLabel?.includes('來因宮')) return '來因宮';
      return '本命';
    default: return '飛星';
  }
}

const PalaceDetailView: React.FC<PalaceDetailViewProps> = ({ palaces, focus, temporalLayers = [], onBack }) => {
  const focusResolved = resolvePalaceEntry(palaces, focus.palaceName);
  const focusPalaceName = normalizePalaceName(focusResolved?.key || focus.palaceName);
  const focusEntry = focusResolved?.entry;
  const focusGan = focusEntry?.palace_gan || '';
  const focusBranch = focus.branch || focusEntry?.branch || '';

  const sanFangPalaces = useMemo(() => {
    return calculateSanFang(focusPalaceName);
  }, [focusPalaceName]);

  // 根據選擇的層級計算飛星四化
  const flyHuaResults = useMemo(() => {
    const results: Array<{type: string; star: string; targetPalace: string; targetBranch: string; source: string; sourceLabel: string}> = [];
    
    // 從 temporalLayers 中提取各層級的 stem
    const dayunLayer = temporalLayers.find(l => l.label === '大限');
    const liunianLayer = temporalLayers.find(l => l.label === '流年');
    const liuyueLayer = temporalLayers.find(l => l.label === '流月');
    const liuriLayer = temporalLayers.find(l => l.label === '流日');
    
    // 判斷當前選擇的層級
    const isNatalLevel = focus.label.includes('來因宮') || focus.label.includes('命宮') || focus.label.includes('身宮');
    const isDaYunLevel = focus.label.includes('大限');
    const isLiuNianLevel = focus.label.includes('流年');
    const isLiuYueLevel = focus.label.includes('流月');
    const isLiuRiLevel = focus.label.includes('流日');
    
    // 本命層級：顯示該宮位的飛星四化
    if (isNatalLevel && focusGan) {
      const hua = calculateFlyHua(focusGan, palaces);
      hua.forEach(h => {
        results.push({
          ...h,
          source: 'natal-fly',
          sourceLabel: '本命'
        });
      });
    }
    
    // 大限層級：顯示大限飛星四化
    if (isDaYunLevel && dayunLayer?.stem) {
      const hua = calculateFlyHua(dayunLayer.stem, palaces);
      hua.forEach(h => {
        results.push({
          ...h,
          source: 'dayun',
          sourceLabel: '大限'
        });
      });
    }
    
    // 流年層級：顯示大限+流年飛星四化
    if (isLiuNianLevel) {
      if (dayunLayer?.stem) {
        const hua = calculateFlyHua(dayunLayer.stem, palaces);
        hua.forEach(h => {
          results.push({
            ...h,
            source: 'dayun',
            sourceLabel: '大限'
          });
        });
      }
      if (liunianLayer?.stem) {
        const hua = calculateFlyHua(liunianLayer.stem, palaces);
        hua.forEach(h => {
          results.push({
            ...h,
            source: 'liunian-fly',
            sourceLabel: '流年'
          });
        });
      }
    }
    
    // 流月層級：顯示大限+流年+流月飛星四化
    if (isLiuYueLevel) {
      if (dayunLayer?.stem) {
        const hua = calculateFlyHua(dayunLayer.stem, palaces);
        hua.forEach(h => {
          results.push({
            ...h,
            source: 'dayun',
            sourceLabel: '大限'
          });
        });
      }
      if (liunianLayer?.stem) {
        const hua = calculateFlyHua(liunianLayer.stem, palaces);
        hua.forEach(h => {
          results.push({
            ...h,
            source: 'liunian-fly',
            sourceLabel: '流年'
          });
        });
      }
      if (liuyueLayer?.stem) {
        const hua = calculateFlyHua(liuyueLayer.stem, palaces);
        hua.forEach(h => {
          results.push({
            ...h,
            source: 'liuyue-fly',
            sourceLabel: '流月'
          });
        });
      }
    }
    
    // 流日層級：顯示全部飛星四化
    if (isLiuRiLevel) {
      if (dayunLayer?.stem) {
        const hua = calculateFlyHua(dayunLayer.stem, palaces);
        hua.forEach(h => {
          results.push({
            ...h,
            source: 'dayun',
            sourceLabel: '大限'
          });
        });
      }
      if (liunianLayer?.stem) {
        const hua = calculateFlyHua(liunianLayer.stem, palaces);
        hua.forEach(h => {
          results.push({
            ...h,
            source: 'liunian-fly',
            sourceLabel: '流年'
          });
        });
      }
      if (liuyueLayer?.stem) {
        const hua = calculateFlyHua(liuyueLayer.stem, palaces);
        hua.forEach(h => {
          results.push({
            ...h,
            source: 'liuyue-fly',
            sourceLabel: '流月'
          });
        });
      }
      if (liuriLayer?.stem) {
        const hua = calculateFlyHua(liuriLayer.stem, palaces);
        hua.forEach(h => {
          results.push({
            ...h,
            source: 'liuri-fly',
            sourceLabel: '流日'
          });
        });
      }
    }
    
    return results;
  }, [focusGan, palaces, focus.label, temporalLayers]);

  const sanFangSet = useMemo(() => new Set(sanFangPalaces), [sanFangPalaces]);
  const inSanFang = flyHuaResults.filter(r => sanFangSet.has(normalizePalaceName(r.targetPalace)));
  const outSanFang = flyHuaResults.filter(r => !sanFangSet.has(normalizePalaceName(r.targetPalace)));

  const outsideNatalTransforms = useMemo(() => {
    const results: Array<{ palaceName: string; branch: string; transform: TransformInfo }> = [];
    Object.entries(palaces).forEach(([palaceName, entry]) => {
      if (!entry || sanFangSet.has(normalizePalaceName(palaceName))) return;
      (entry.natal_transforms || []).forEach((t) => {
        results.push({
          palaceName: normalizePalaceName(palaceName),
          branch: entry.branch,
          transform: t,
        });
      });
    });
    return results;
  }, [palaces, sanFangSet]);

  // 收集非三方四正的流運四化（從後端返回的 dayun_transforms 等字段）
  const outsideTemporalTransforms = useMemo(() => {
    const results: Array<{ palaceName: string; branch: string; transform: TransformInfo; source: string }> = [];
    
    Object.entries(palaces).forEach(([palaceName, entry]) => {
      if (!entry || sanFangSet.has(normalizePalaceName(palaceName))) return;
      
      // 大限四化
      (entry.dayun_transforms || []).forEach((t) => {
        results.push({
          palaceName: normalizePalaceName(palaceName),
          branch: entry.branch,
          transform: t,
          source: 'dayun'
        });
      });
      
      // 流年四化
      (entry.liunian_transforms || []).forEach((t) => {
        results.push({
          palaceName: normalizePalaceName(palaceName),
          branch: entry.branch,
          transform: t,
          source: 'liunian'
        });
      });
      
      // 流月四化
      (entry.liuyue_transforms || []).forEach((t) => {
        results.push({
          palaceName: normalizePalaceName(palaceName),
          branch: entry.branch,
          transform: t,
          source: 'liuyue'
        });
      });
      
      // 流日四化
      (entry.liuri_transforms || []).forEach((t) => {
        results.push({
          palaceName: normalizePalaceName(palaceName),
          branch: entry.branch,
          transform: t,
          source: 'liuri'
        });
      });
    });
    return results;
  }, [palaces, sanFangSet]);

  const temporalNarratives = useMemo(() => {
    const temporalFocusLabels = new Set(['大限', '流年', '流月', '流日']);
    const scopedLayers = temporalFocusLabels.has(focus.label)
      ? temporalLayers.filter((layer) => layer.label === focus.label)
      : temporalLayers;

    return scopedLayers
      .filter((layer) => layer.palace)
      .map((layer) => {
        const role = getRoleLabel(layer.palace as string, focusPalaceName, sanFangPalaces);
        const mappedPalace = resolveShiftedPalace(layer.palace as string, focusPalaceName);
        const roleText = role || '三方四正外宮';
        const baseTime = layer.timeLabel ? `（${layer.timeLabel}）` : '';
        const stemText = layer.stem ? `${layer.stem}干` : '未知干';
        const timeGanzhi = layer.stem && layer.timeBranch ? `${layer.stem}${layer.timeBranch}` : '未知';
        const landingText = layer.branch ? `${layer.branch}宮` : '未知宮';
        const flyStarHint = role
          ? `${layer.label}${baseTime}落在${roleText}，以${stemText}啟動飛星四化，優先觀察該宮與本宮互動。`
          : `${layer.label}${baseTime}落在三方四正外宮，飛星影響偏向外部情境，需搭配外宮飛入訊號判讀。`;
        const sanHeHint = role
          ? `${layer.label}落入${roleText}，三合派可視為本局同盟位活化，重點看該宮主題如何回流本宮。`
          : `${layer.label}不在本宮/三合/對宮，三合派以旁證參考，不作主論斷。`;

        return {
          key: `${layer.label}-${layer.palace}-${layer.branch}`,
          title: `${layer.label}${baseTime} → ${layer.palace}（${layer.branch || '未知'}）｜→ ${mappedPalace || '—'}`,
          temporalLine: `時間干支：${timeGanzhi}｜落宮：${landingText}`,
          flyStarHint,
          sanHeHint,
        };
      });
  }, [temporalLayers, focus.label, focusPalaceName, sanFangPalaces]);

  const getPalaceStars = (palaceName: string): StarDetail[] => {
    const entry = resolvePalaceEntry(palaces, palaceName)?.entry;
    if (!entry) return [];

    const stars: StarDetail[] = [];

    if (entry.star_details && entry.star_details.length > 0) {
      entry.star_details.forEach(star => {
        stars.push({
          name: star.name,
          brightness: star.brightness,
          isMain: mainStarSet.has(star.name)
        });
      });
    } else if (entry.stars) {
      entry.stars.forEach(star => {
        stars.push({
          name: star,
          isMain: mainStarSet.has(star)
        });
      });
    }

    if (entry.assistant_stars) {
      entry.assistant_stars.forEach(star => {
        stars.push({
          name: star,
          isMain: false
        });
      });
    }

    if (entry.secondary_stars) {
      entry.secondary_stars.forEach(star => {
        stars.push({
          name: star,
          isMain: false
        });
      });
    }

    return stars;
  };

  const getPalaceTransforms = (palaceName: string): TransformInfo[] => {
    const entry = resolvePalaceEntry(palaces, palaceName)?.entry;
    if (!entry) return [];

    const transforms: TransformInfo[] = [];

    if (entry.natal_transforms) {
      entry.natal_transforms.forEach(t => {
        transforms.push({
          star: t.star,
          transformation: t.transformation,
          display: t.display,
          source: 'natal'
        });
      });
    }

    if (focus.label.includes('命宮') || focus.label.includes('身宮') || focus.label.includes('來因宮')) {
      const palaceEntry = resolvePalaceEntry(palaces, focus.palaceName)?.entry;
      if (palaceEntry?.palace_gan) {
        const flyResults = calculateFlyHua(palaceEntry.palace_gan, palaces);
        flyResults.forEach(fly => {
          const alreadyExists = transforms.some(t => 
            t.star === fly.star && t.transformation === '化' + fly.type
          );
          if (!alreadyExists) {
            transforms.push({
              star: fly.star,
              transformation: '化' + fly.type,
              display: fly.star + '化' + fly.type,
              source: 'natal-fly'
            });
          }
        });
      }
    }

    if (entry.dayun_transforms) {
      entry.dayun_transforms.forEach(t => {
        transforms.push({
          star: t.star,
          transformation: t.transformation,
          display: t.display,
          source: 'dayun'
        });
      });
    }

    if (entry.liunian_transforms) {
      entry.liunian_transforms.forEach(t => {
        transforms.push({
          star: t.star,
          transformation: t.transformation,
          display: t.display,
          source: 'liunian'
        });
      });
    }

    if (entry.liuyue_transforms) {
      entry.liuyue_transforms.forEach(t => {
        transforms.push({
          star: t.star,
          transformation: t.transformation,
          display: t.display,
          source: 'liuyue'
        });
      });
    }

    if (entry.liuri_transforms) {
      entry.liuri_transforms.forEach(t => {
        transforms.push({
          star: t.star,
          transformation: t.transformation,
          display: t.display,
          source: 'liuri'
        });
      });
    }

    return transforms;
  };

  const getStarTransforms = (starName: string, palaceName: string): TransformInfo[] => {
    const palaceTransforms = getPalaceTransforms(palaceName);
    return palaceTransforms.filter(t => t.star === starName);
  };

  const getFlyInTransforms = (palaceName: string): Array<{type: string; star: string; source: string; sourceLabel: string}> => {
    const normalizedPalace = normalizePalaceName(palaceName);
    const flyIn = inSanFang.filter(f => normalizePalaceName(f.targetPalace) === normalizedPalace);
    return flyIn.map(f => ({
      type: f.type,
      star: f.star,
      source: f.source,
      sourceLabel: f.sourceLabel
    }));
  };

  const getTitleExtra = (): string => {
    const parts: string[] = [];
    if (focus.temporalInfo) {
      if (focus.temporalInfo.stem && focus.temporalInfo.timeBranch) {
        parts.push(`${focus.temporalInfo.stem}${focus.temporalInfo.timeBranch}`);
      } else if (focus.temporalInfo.stem) {
        parts.push(`${focus.temporalInfo.stem}干`);
      }
      if (focus.temporalInfo.timeLabel) {
        parts.push(focus.temporalInfo.timeLabel);
      }
      if (focus.temporalInfo.ageRange) {
        parts.push(`${focus.temporalInfo.ageRange}`);
      }
    }
    return parts.length > 0 ? ` ${parts.join(' · ')}` : '';
  };

  const getBrightnessColor = (brightness?: string): string => {
    switch (brightness) {
      case '廟': return '#16a34a';
      case '旺': return '#22c55e';
      case '得': return '#3b82f6';
      case '利': return '#6366f1';
      case '平': return '#6b7280';
      case '陷': return '#dc2626';
      case '不': return '#9ca3af';
      default: return '#6b7280';
    }
  };

  const showTemporalNarratives = !focus.label.startsWith('來因宮') && !focus.label.startsWith('命宮') && !focus.label.startsWith('身宮');

  return (
    <div className="section-stack">
      <button
        className="btn-secondary"
        onClick={onBack}
        style={{
          background: 'var(--primary)',
          color: 'var(--cream)',
          border: 'none',
          padding: '0.5rem 1rem',
          borderRadius: 'var(--radius-sm)',
          display: 'flex',
          alignItems: 'center',
          gap: '0.5rem',
          cursor: 'pointer',
          fontWeight: 500,
        }}
      >
        <ChevronLeft size={16} />
        返回命盤總覽
      </button>

      <div className="heading-md">
        {focus.label} · {focusGan || '—'}{focusBranch || '—'}
        <span className="body-sm" style={{ color: 'var(--text-soft)', marginLeft: '0.5rem' }}>
          {getTitleExtra()}
        </span>
        {focusGan && (
          <span className="body-sm" style={{ color: 'var(--text-soft)', marginLeft: '0.5rem' }}>
            （{focusGan}干四化飛布）
          </span>
        )}
      </div>

      <div
        style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(4, minmax(220px, 1fr))',
          gap: '1rem',
          overflowX: 'auto',
        }}
      >
        {sanFangPalaces.map((palaceName) => {
          const resolved = resolvePalaceEntry(palaces, palaceName);
          if (!resolved) return null;
          const { entry } = resolved;

          const stars = getPalaceStars(palaceName);
          const flyInTransforms = getFlyInTransforms(palaceName);
          const roleLabel = getRoleLabel(palaceName, focusPalaceName, sanFangPalaces);
          const mappedPalace = resolveShiftedPalace(palaceName, focusPalaceName);

          return (
            <div key={palaceName} className="metric-card" style={{ border: palaceName === focusPalaceName ? '2px solid var(--cta)' : undefined }}>
              <div className="metric-label">{roleLabel} · {palaceName}（{entry.branch}）</div>
              <div className="body-sm" style={{ color: 'var(--text-soft)', marginTop: '0.2rem' }}>
                → {mappedPalace || '—'}
              </div>
              <div style={{ marginTop: '0.75rem', display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
                {stars.map((star, index) => {
                  const starTransforms = getStarTransforms(star.name, palaceName);
                  const starFlyIn = flyInTransforms.filter(f => f.star === star.name);

                  return (
                    <div
                      key={`${palaceName}-${star.name}-${index}`}
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.5rem',
                        flexWrap: 'wrap',
                        padding: '0.35rem 0',
                        borderBottom: index < stars.length - 1 ? '1px solid rgba(178, 135, 70, 0.1)' : 'none'
                      }}
                    >
                      <span
                        className="body-md"
                        style={{
                          fontWeight: star.isMain ? 600 : 400,
                          color: star.isMain ? 'var(--primary)' : 'var(--secondary)',
                          minWidth: '2.5rem'
                        }}
                      >
                        {star.name}
                      </span>

                      {star.brightness && star.isMain && (
                        <span
                          className="body-sm"
                          style={{
                            color: getBrightnessColor(star.brightness),
                            fontWeight: 500,
                            minWidth: '1.5rem'
                          }}
                        >
                          {star.brightness}
                        </span>
                      )}

                      {starTransforms.length > 0 && (
                        <div style={{ display: 'flex', gap: '0.25rem', flexWrap: 'wrap' }}>
                          {starTransforms.map((t, i) => (
                            <span
                              key={i}
                              className="body-sm"
                              style={{
                                padding: '0.1rem 0.35rem',
                                borderRadius: '4px',
                                fontWeight: 500,
                                background: t.source === 'natal' ? 'rgba(220, 38, 38, 0.1)' :
                                           t.source === 'dayun' ? 'rgba(139, 92, 246, 0.1)' :
                                           t.source === 'liunian' ? 'rgba(59, 130, 246, 0.1)' :
                                           t.source === 'liuyue' ? 'rgba(16, 185, 129, 0.1)' :
                                           t.source === 'liuri' ? 'rgba(245, 158, 11, 0.1)' :
                                           'rgba(37, 99, 235, 0.1)',
                                color: t.source === 'natal' ? '#dc2626' :
                                       t.source === 'dayun' ? '#8b5cf6' :
                                       t.source === 'liunian' ? '#2563eb' :
                                       t.source === 'liuyue' ? '#10b981' :
                                       t.source === 'liuri' ? '#f59e0b' :
                                       '#2563eb'
                              }}
                            >
                              {getTransformSourceLabel(t.source, focus.label)}{t.transformation.replace('化', '')}
                            </span>
                          ))}
                        </div>
                      )}

                      {starFlyIn.map((f, i) => (
                        <span
                          key={`fly-${i}`}
                          className="body-sm"
                          style={{
                            padding: '0.1rem 0.35rem',
                            borderRadius: '4px',
                            fontWeight: 500,
                            background: f.source === 'dayun' ? 'rgba(139, 92, 246, 0.1)' :
                                       f.source === 'liunian-fly' ? 'rgba(59, 130, 246, 0.1)' :
                                       f.source === 'liuyue-fly' ? 'rgba(16, 185, 129, 0.1)' :
                                       f.source === 'liuri-fly' ? 'rgba(245, 158, 11, 0.1)' :
                                       'rgba(37, 99, 235, 0.1)',
                            color: f.source === 'dayun' ? '#8b5cf6' :
                                   f.source === 'liunian-fly' ? '#2563eb' :
                                   f.source === 'liuyue-fly' ? '#10b981' :
                                   f.source === 'liuri-fly' ? '#f59e0b' :
                                   '#2563eb'
                          }}
                        >
                          {f.sourceLabel}{f.type}
                        </span>
                      ))}
                    </div>
                  );
                })}

                {flyInTransforms.filter(f => !stars.some(s => s.name === f.star)).map((f, i) => {
                  const flyColor = f.source === 'dayun' ? '#8b5cf6' :
                                    f.source === 'liunian-fly' ? '#2563eb' :
                                    f.source === 'liuyue-fly' ? '#10b981' :
                                    f.source === 'liuri-fly' ? '#f59e0b' :
                                    '#2563eb';
                  return (
                    <div
                      key={`fly-extra-${i}`}
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.5rem',
                        padding: '0.35rem 0',
                        color: flyColor,
                        fontSize: '0.875rem'
                      }}
                    >
                      <span>飛入：</span>
                      <span style={{ fontWeight: 500 }}>{f.sourceLabel}{f.type}</span>
                    </div>
                  );
                })}
              </div>
            </div>
          );
        })}
      </div>

      {showTemporalNarratives && temporalNarratives.length > 0 && (
        <div className="card section-stack">
          <div className="heading-sm">流運落宮補充說明（飛星主述 / 三合補充）</div>
          {temporalNarratives.map((item) => (
            <div key={item.key} className="metric-card" style={{ display: 'grid', gap: '0.35rem' }}>
              <div className="metric-label">{item.title}</div>
              <div className="body-sm" style={{ color: 'var(--text-soft)' }}>{item.temporalLine}</div>
              <div className="body-sm" style={{ color: '#2563eb' }}>飛星派：{item.flyStarHint}</div>
              <div className="body-sm" style={{ color: '#7c3aed' }}>三合派：{item.sanHeHint}</div>
            </div>
          ))}
        </div>
      )}

      {(outsideNatalTransforms.length > 0 || outSanFang.length > 0) && (
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
            gap: '1rem',
            alignItems: 'stretch',
          }}
        >
          {/* 生年四化（不在三方四正） */}
          {outsideNatalTransforms.length > 0 && (
            <div className="card" style={{ display: 'flex', flexDirection: 'column', height: '100%', gap: '0.5rem' }}>
              <div className="heading-sm" style={{ marginBottom: '0.5rem' }}>非三方四正的生年四化</div>
              <div style={{ display: 'grid', gap: '0.5rem', alignContent: 'start', flex: 1 }}>
                {outsideNatalTransforms.map((item, index) => {
                  const stars = getPalaceStars(item.palaceName);
                  const starDetail = stars.find(s => s.name === item.transform.star);
                  const brightness = starDetail?.brightness;
                  const shiftedPalace = resolveShiftedPalace(item.palaceName, focusPalaceName);
                  return (
                    <div key={`${item.palaceName}-${index}`} className="metric-card">
                      <div className="body-md">
                        {item.palaceName}（{item.branch}）→ {shiftedPalace || '—'}
                      </div>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginTop: '0.25rem' }}>
                        <span style={{ fontWeight: 600, color: 'var(--primary)' }}>{item.transform.star}</span>
                        {brightness && (
                          <span style={{ color: getBrightnessColor(brightness), fontWeight: 500 }}>{brightness}</span>
                        )}
                        <span style={{ color: '#dc2626', fontWeight: 500 }}>生年{item.transform.transformation.replace('化', '')}</span>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          )}

          {(outsideTemporalTransforms.length > 0 || outSanFang.length > 0) && (
            <div className="card" style={{ display: 'flex', flexDirection: 'column', height: '100%', gap: '0.5rem' }}>
              <div className="heading-sm" style={{ marginBottom: '0.5rem' }}>
                非三方四正的
                {focus.label.includes('大限') ? '大限四化' : 
                 focus.label.includes('流年') ? '流年四化' : 
                 focus.label.includes('流月') ? '流月四化' : 
                 focus.label.includes('流日') ? '流日四化' : 
                 focus.label.includes('命宮') ? '命宮四化' : 
                 focus.label.includes('身宮') ? '身宮四化' : 
                 focus.label.includes('來因宮') ? '來因宮四化' : '四化'}
              </div>
              <div style={{ display: 'grid', gap: '0.5rem', alignContent: 'start', flex: 1 }}>
                {outsideTemporalTransforms.map((item, index) => {
                  const stars = getPalaceStars(item.palaceName);
                  const starDetail = stars.find(s => s.name === item.transform.star);
                  const brightness = starDetail?.brightness;
                  const shiftedPalace = resolveShiftedPalace(item.palaceName, focusPalaceName);
                  const sourceLabel = item.source === 'dayun' ? '大限' : 
                                     item.source === 'liunian' ? '流年' : 
                                     item.source === 'liuyue' ? '流月' : 
                                     item.source === 'liuri' ? '流日' : '';
                  const sourceColor = item.source === 'dayun' ? '#8b5cf6' : 
                                     item.source === 'liunian' ? '#2563eb' : 
                                     item.source === 'liuyue' ? '#10b981' : 
                                     item.source === 'liuri' ? '#f59e0b' : '#2563eb';
                  return (
                    <div key={`temporal-${item.palaceName}-${index}`} className="metric-card">
                      <div className="body-md">
                        {item.palaceName}（{item.branch}）→ {shiftedPalace || '—'}
                      </div>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginTop: '0.25rem' }}>
                        <span style={{ fontWeight: 600, color: 'var(--primary)' }}>{item.transform.star}</span>
                        {brightness && (
                          <span style={{ color: getBrightnessColor(brightness), fontWeight: 500 }}>{brightness}</span>
                        )}
                        <span style={{ color: sourceColor, fontWeight: 500 }}>
                          {sourceLabel}{item.transform.transformation.replace('化', '')}
                        </span>
                      </div>
                    </div>
                  );
                })}
                
                {outSanFang.map((f, i) => {
                  const stars = getPalaceStars(f.targetPalace);
                  const starDetail = stars.find(s => s.name === f.star);
                  const brightness = starDetail?.brightness;
                  const shiftedPalace = resolveShiftedPalace(f.targetPalace, focusPalaceName);
                  const sourceLabel = f.source === 'dayun' ? '大限' : 
                                     f.source === 'liunian-fly' ? '流年' : 
                                     f.source === 'liuyue-fly' ? '流月' : 
                                     f.source === 'liuri-fly' ? '流日' : 
                                     f.source === 'natal-fly' ? (focus.label.includes('命宮') ? '命宮' : focus.label.includes('身宮') ? '身宮' : focus.label.includes('來因宮') ? '來因宮' : '本命') : '';
                  const flyColor = f.source === 'dayun' ? '#8b5cf6' : 
                                    f.source === 'liunian-fly' ? '#2563eb' : 
                                    f.source === 'liuyue-fly' ? '#10b981' : 
                                    f.source === 'liuri-fly' ? '#f59e0b' : 
                                    f.source === 'natal-fly' ? '#06b6d4' : '#2563eb';
                  return (
                    <div key={`fly-${f.star}-${f.type}-${i}`} className="metric-card">
                      <div className="body-md">
                        {f.targetPalace}（{f.targetBranch}）→ {shiftedPalace || '—'}
                      </div>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginTop: '0.25rem' }}>
                        <span style={{ fontWeight: 600, color: 'var(--primary)' }}>{f.star}</span>
                        {brightness && (
                          <span style={{ color: getBrightnessColor(brightness), fontWeight: 500 }}>{brightness}</span>
                        )}
                        <span style={{ color: flyColor, fontWeight: 500 }}>
                          {sourceLabel}{f.type}
                        </span>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          )}
        </div>
      )}

      <div className="body-sm" style={{ display: 'flex', gap: '1rem', color: 'var(--secondary)', flexWrap: 'wrap' }}>
        <span><span style={{ color: '#16a34a' }}>●</span> 廟</span>
        <span><span style={{ color: '#22c55e' }}>●</span> 旺</span>
        <span><span style={{ color: '#3b82f6' }}>●</span> 得</span>
        <span><span style={{ color: '#6366f1' }}>●</span> 利</span>
        <span><span style={{ color: '#6b7280' }}>●</span> 平</span>
        <span><span style={{ color: '#dc2626' }}>●</span> 陷</span>
        <span style={{ marginLeft: '1rem' }}>|</span>
        <span><span style={{ color: '#dc2626' }}>■</span> 生年四化（祿權科忌）</span>
        {focus.label.includes('大限') || focus.label.includes('流年') || focus.label.includes('流月') || focus.label.includes('流日') ? (
          <span><span style={{ color: '#8b5cf6' }}>■</span> 大限四化</span>
        ) : null}
        {focus.label.includes('流年') || focus.label.includes('流月') || focus.label.includes('流日') ? (
          <span><span style={{ color: '#2563eb' }}>■</span> 流年四化</span>
        ) : null}
        {focus.label.includes('流月') || focus.label.includes('流日') ? (
          <span><span style={{ color: '#10b981' }}>■</span> 流月四化</span>
        ) : null}
        {focus.label.includes('流日') ? (
          <span><span style={{ color: '#f59e0b' }}>■</span> 流日四化</span>
        ) : null}
        {focus.label.includes('來因宮') || focus.label.includes('命宮') || focus.label.includes('身宮') ? (
          <span><span style={{ color: '#06b6d4' }}>■</span> 
            {focus.label.includes('命宮') ? '命宮' : 
             focus.label.includes('身宮') ? '身宮' : 
             focus.label.includes('來因宮') ? '來因宮' : '本命'}四化（{focusGan}干）
          </span>
        ) : null}
      </div>
    </div>
  );
};

export default PalaceDetailView;
