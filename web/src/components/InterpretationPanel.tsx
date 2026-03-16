import React from 'react';
import { motion } from 'framer-motion';
import { Sparkles, ArrowRight, Star, Crown, Compass } from 'lucide-react';
import '../styles/design-system.css';

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

interface InterpretationData {
  summary: string;
  character_traits: string;
  origin_palace_analysis: string;
  karmic_narrative: KarmicStep[];
  san_fang_diagnosis: SanFangRole[];
  star_details?: DeepStarAnalysis[];
  origin_fly_hua?: any;
  temporal_resonance: any[];
  classic_patterns?: string[];
}

interface InterpretationPanelProps {
  interpretation: InterpretationData;
  className?: string;
}

const InterpretationPanel: React.FC<InterpretationPanelProps> = ({
  interpretation,
  className = ''
}) => {
  const getTypeColor = (type: string) => {
    const colors: Record<string, string> = {
      '祿': 'text-yellow-600 bg-yellow-50',
      '權': 'text-blue-600 bg-blue-50',
      '科': 'text-green-600 bg-green-50',
      '忌': 'text-red-600 bg-red-50'
    };
    return colors[type] || 'text-slate-600 bg-slate-50';
  };

  const getTypeIcon = (type: string) => {
    switch (type) {
      case '祿':
        return <Crown className="w-4 h-4" />;
      case '權':
        return <Compass className="w-4 h-4" />;
      case '科':
        return <Star className="w-4 h-4" />;
      case '忌':
        return <ArrowRight className="w-4 h-4" />;
      default:
        return <Sparkles className="w-4 h-4" />;
    }
  };

  return (
    <div className={`interpretation-panel ${className}`}>
      {/* 大師總論 */}
      <motion.section 
        className="card card-gold mb-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
      >
        <div className="flex items-center gap-3 mb-4">
          <Sparkles className="w-6 h-6 text-yellow-600" />
          <h2 className="heading-md">大師總論</h2>
        </div>
        <p className="body-lg leading-relaxed">{interpretation.summary}</p>
      </motion.section>

      {/* 來因宮分析 */}
      <motion.section 
        className="card card-navy mb-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6, delay: 0.1 }}
      >
        <div className="flex items-center gap-3 mb-4">
          <Crown className="w-6 h-6 text-amber-700" />
          <h2 className="heading-md">來因宮動態因果</h2>
        </div>
        <p className="body-md leading-relaxed">
          {interpretation.origin_palace_analysis}
        </p>
      </motion.section>

      {/* 能量循環 (祿隨忌走) */}
      <motion.section 
        className="mb-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6, delay: 0.2 }}
      >
        <h2 className="heading-md mb-6 flex items-center gap-3">
          <ArrowRight className="w-6 h-6 text-blue-600" />
          能量循環 (祿隨忌走)
        </h2>
        <div className="grid gap-4">
          {interpretation.karmic_narrative.map((step, index) => (
            <motion.div
              key={index}
              className="card border-l-4 border-blue-600"
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.4, delay: 0.3 + index * 0.1 }}
            >
              <div className="flex items-start gap-4">
                <div className={`p-2 rounded-full ${getTypeColor(step.type)}`}>
                  {getTypeIcon(step.type)}
                </div>
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2">
                    <span className={`px-2 py-1 rounded-full text-xs font-bold ${getTypeColor(step.type)}`}>
                      {step.type}
                    </span>
                    <span className="text-sm font-semibold text-slate">
                      {step.role}
                    </span>
                  </div>
                  <h3 className="font-bold text-primary mb-1">
                    {step.star} 落入 {step.palace}
                  </h3>
                  <p className="body-sm text-slate leading-relaxed">
                    {step.desc}
                  </p>
                </div>
              </div>
            </motion.div>
          ))}
        </div>
      </motion.section>

      {/* 三方四正專業診斷 */}
      <motion.section 
        className="mb-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6, delay: 0.4 }}
      >
        <h2 className="heading-md mb-6 flex items-center gap-3">
          <Compass className="w-6 h-6 text-blue-600" />
          三方四正專業診斷
        </h2>
        <div className="grid md:grid-cols-2 gap-4">
          {interpretation.san_fang_diagnosis.map((role, index) => (
            <motion.div
              key={index}
              className="card bg-gradient-to-br from-slate-50 to-white"
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ duration: 0.4, delay: 0.5 + index * 0.1 }}
            >
              <div className="flex items-start gap-3">
                <div className="w-2 h-2 bg-blue-600 rounded-full mt-2"></div>
                <div>
                  <h3 className="font-bold text-primary mb-1">
                    {role.role}
                  </h3>
                  <p className="text-sm text-slate mb-2">
                    【{role.palace}】
                  </p>
                  <p className="body-sm text-slate leading-relaxed">
                    {role.diagnosis}
                  </p>
                </div>
              </div>
            </motion.div>
          ))}
        </div>
      </motion.section>

      {/* 星曜深度解析 */}
      {interpretation.star_details && interpretation.star_details.length > 0 && (
        <motion.section 
          className="mb-8"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.6 }}
        >
          <h2 className="heading-md mb-6 flex items-center gap-3">
            <Star className="w-6 h-6 text-yellow-600" />
            星曜深度解析
          </h2>
          <div className="space-y-6">
            {interpretation.star_details.map((star, index) => (
              <motion.div
                key={index}
                className="card card-gold"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.4, delay: 0.7 + index * 0.1 }}
              >
                <h3 className="heading-sm mb-3 text-center">
                  【 {star.name} 】
                </h3>
                
                {star.verse && (
                  <div className="mb-4 p-3 bg-cream rounded-lg border-l-4 border-yellow-600">
                    <p className="body-sm italic text-slate text-center">
                      {star.verse}
                    </p>
                  </div>
                )}
                
                <div className="grid md:grid-cols-3 gap-4">
                  <div className="text-center">
                    <h4 className="font-bold text-green-700 mb-2">正面特質</h4>
                    <p className="body-sm text-slate">{star.positive}</p>
                  </div>
                  <div className="text-center">
                    <h4 className="font-bold text-red-700 mb-2">負面特質</h4>
                    <p className="body-sm text-slate">{star.negative}</p>
                  </div>
                  <div className="text-center">
                    <h4 className="font-bold text-blue-700 mb-2">修行之道</h4>
                    <p className="body-sm text-slate">{star.remedy}</p>
                  </div>
                </div>
                
                {star.evolution && (
                  <div className="mt-4 p-3 bg-blue-50 rounded-lg">
                    <p className="body-sm text-blue-800">{star.evolution}</p>
                  </div>
                )}
              </motion.div>
            ))}
          </div>
        </motion.section>
      )}

      {/* 經典格局 */}
      {interpretation.classic_patterns && interpretation.classic_patterns.length > 0 && (
        <motion.section 
          className="mb-8"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.8 }}
        >
          <h2 className="heading-md mb-6 flex items-center gap-3">
            <Crown className="w-6 h-6 text-yellow-600" />
            檢出格局
          </h2>
          <div className="flex flex-wrap gap-3">
            {interpretation.classic_patterns.map((pattern, index) => (
              <motion.div
                key={index}
                className="px-4 py-2 bg-gradient-to-r from-yellow-100 to-yellow-50 border border-yellow-300 rounded-full"
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.3, delay: 0.9 + index * 0.05 }}
              >
                <span className="text-sm font-bold text-yellow-800">{pattern}</span>
              </motion.div>
            ))}
          </div>
        </motion.section>
      )}
    </div>
  );
};

export default InterpretationPanel;
