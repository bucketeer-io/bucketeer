import { FC, memo, useCallback, useState } from 'react';

import { classNames } from '../../utils/css';
import { CopyChip } from '../CopyChip';

export interface TagChipsProps {
  tags: string[];
}

export const TagChips: FC<TagChipsProps> = memo(({ tags }) => {
  return (
    <div className="flex flex-row flex-wrap space-x-1">
      {tags.map((tag) => {
        return (
          <CopyChip key={tag} text={tag}>
            <div
              key={`${tag}`}
              className={classNames(
                'text-xs relative inline-block p-1',
                'leading-tight cursor-pointer'
              )}
            >
              <span
                aria-hidden="true"
                className={classNames(
                  'absolute inset-0',
                  'border border-gray-300 rounded'
                )}
              ></span>
              <span className="relative text-gray-700">{tag}</span>
            </div>
          </CopyChip>
        );
      })}
    </div>
  );
});
