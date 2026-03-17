import { useTranslation } from 'i18n';
import { EvaluationFeature } from 'pages/debugger/types';
import { FlagVariationPolygon } from 'pages/feature-flags/collection-layout/elements';
import { Card } from 'components/mobile-card/elements';
import ReasonTooltip from './reason-tooltip';
import ResultName from './result-name';

interface Props {
  evaluation: EvaluationFeature;
  index: number;
  isFlag: boolean;
  handleGetMaintainerInfo: (email: string) => string;
}

const ResultItemCard = ({
  evaluation,
  index,
  isFlag,
  handleGetMaintainerInfo
}: Props) => {
  const { t } = useTranslation(['common', 'table', 'form']);
  const {
    featureId,
    userId,
    feature,
    variationName,
    variationValue,
    variationId,
    reason
  } = evaluation;

  return (
    <Card className="shadow-none border border-gray-200 rounded-xl p-4">
      <Card.Header
        triger={
          <ResultName
            feature={feature}
            id={!isFlag ? featureId : userId}
            isFlag={isFlag}
            name={!isFlag ? feature.name : userId}
            variationType={feature.variationType}
            maintainer={handleGetMaintainerInfo(feature.maintainer)}
            onTable
          />
        }
      />
      <Card.Meta>
        <div className="flex flex-col gap-y-3">
          <div className="flex flex-col gap-y-1">
            <p className="typo-para-tiny font-bold uppercase text-gray-500">
              {t('table:feature-flags.variation')}
            </p>
            <div className="flex items-start gap-x-1.5">
              <FlagVariationPolygon index={index} className="mt-0.5 shrink-0" />
              <div className="flex flex-col gap-y-1">
                <p className="typo-para-small text-gray-700 break-all">
                  {variationName || variationValue}
                </p>
                <p className="typo-para-small text-gray-500 break-all">
                  {variationId}
                </p>
              </div>
            </div>
          </div>
          <div className="flex flex-col gap-y-1">
            <p className="typo-para-tiny font-bold uppercase text-gray-500">
              {t('table:reason')}
            </p>
            <ReasonTooltip reason={reason} />
          </div>
        </div>
      </Card.Meta>
    </Card>
  );
};

export default ResultItemCard;
