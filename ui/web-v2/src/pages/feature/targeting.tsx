import { createVariationLabel } from '../../utils/variation';
import { yupResolver } from '@hookform/resolvers/yup';
import deepEqual from 'deep-equal';
import React, { useCallback, useState, FC, memo, useEffect } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { useIntl } from 'react-intl';
import { useDispatch, useSelector } from 'react-redux';

import {
  FeatureConfirmDialog,
  SaveFeatureType
} from '../../components/FeatureConfirmDialog';
import { FeatureTargetingForm } from '../../components/FeatureTargetingForm';
import { intl } from '../../lang';
import { messages } from '../../lang/messages';
import { AppState } from '../../modules';
import {
  selectById as selectFeatureById,
  updateFeatureTargeting,
  getFeature,
  createCommand,
  updateFeature
} from '../../modules/features';
import { useCurrentEnvironment } from '../../modules/me';
import { listSegments } from '../../modules/segments';
import { Clause } from '../../proto/feature/clause_pb';
import {
  AddClauseCommand,
  AddClauseValueCommand,
  AddPrerequisiteCommand,
  AddRuleCommand,
  AddUserToVariationCommand,
  ChangeClauseAttributeCommand,
  ChangeClauseOperatorCommand,
  ChangeDefaultStrategyCommand,
  ChangeFixedStrategyCommand,
  ChangeOffVariationCommand,
  ChangeRolloutStrategyCommand,
  ChangeRuleStrategyCommand,
  Command,
  DeleteClauseCommand,
  DeleteRuleCommand,
  DisableFeatureCommand,
  EnableFeatureCommand,
  RemoveClauseValueCommand,
  RemoveUserFromVariationCommand,
  ResetSamplingSeedCommand,
  RemovePrerequisiteCommand,
  ChangePrerequisiteVariationCommand,
  ChangeRulesOrderCommand
} from '../../proto/feature/command_pb';
import { Feature } from '../../proto/feature/feature_pb';
import { Prerequisite } from '../../proto/feature/prerequisite_pb';
import { Rule } from '../../proto/feature/rule_pb';
import {
  FixedStrategy,
  RolloutStrategy,
  Strategy
} from '../../proto/feature/strategy_pb';
import { Variation } from '../../proto/feature/variation_pb';
import { AppDispatch } from '../../store';

import {
  ruleClauseType,
  targetingFormSchema,
  RuleClauseType,
  StrategySchema,
  RuleClauseSchema,
  RuleSchema,
  TargetingForm
} from './formSchema';
import {
  ChangeType,
  PrerequisiteChange,
  RuleChange,
  TargetChange
} from '../../proto/feature/service_pb';
import { Target } from '../../proto/feature/target_pb';

interface FeatureTargetingPageProps {
  featureId: string;
}

export const FeatureTargetingPage: FC<FeatureTargetingPageProps> = memo(
  ({ featureId }) => {
    const { formatMessage: f } = useIntl();
    const dispatch = useDispatch<AppDispatch>();

    const currentEnvironment = useCurrentEnvironment();
    const feature = useSelector<AppState, Feature.AsObject | undefined>(
      (state) => {
        const f = selectFeatureById(state.features, featureId);
        return f;
      },
      (left: Feature.AsObject | undefined, right): boolean => {
        return JSON.stringify(left) === JSON.stringify(right);
      }
    );

    const getDefaultValues = (
      feature: Feature.AsObject,
      requireComment: boolean
    ): TargetingForm => {
      return {
        prerequisites: [
          ...new Map(
            feature.prerequisitesList.map((p) => [p.featureId, p])
          ).values()
        ], // remove duplicate prerequisites
        enabled: feature.enabled,
        targets: feature.targetsList.map((t) => {
          return {
            variationId: t.variation,
            users: t.usersList
          };
        }),
        rules: feature.rulesList.map((r) => {
          return {
            id: r.id,
            strategy: createStrategyDefaultValue(
              r.strategy,
              feature.variationsList
            ),
            clauses: r.clausesList.map((c) => {
              return {
                id: c.id,
                type: createClauseType(c.operator),
                attribute: c.attribute,
                operator: c.operator.toString(),
                values: c.valuesList
              };
            })
          };
        }),
        defaultStrategy: createStrategyDefaultValue(
          feature.defaultStrategy,
          feature.variationsList
        ),
        offVariation: feature.offVariation && {
          value: feature.offVariation,
          label: createVariationLabel(
            feature.variationsList.find((v) => v.id === feature.offVariation)
          )
        },
        requireComment: requireComment,
        comment: ''
      };
    };

    const methods = useForm({
      resolver: yupResolver(targetingFormSchema),
      defaultValues: getDefaultValues(
        feature,
        currentEnvironment.requireComment
      ),
      mode: 'onChange'
    });
    const {
      handleSubmit,
      formState: { dirtyFields },
      reset
    } = methods;
    const [isConfirmDialogOpen, setIsConfirmDialogOpen] = useState(false);

    const handleUpdate = useCallback(
      async (data, saveFeatureType) => {
        const prepareUpdate = async (actionType, payload) => {
          dispatch(actionType(payload)).then(() => {
            setIsConfirmDialogOpen(false);
            dispatch(
              getFeature({
                environmentId: currentEnvironment.id,
                id: featureId
              })
            ).then((response) => {
              const featurePayload = response.payload as Feature.AsObject;
              reset(
                getDefaultValues(
                  featurePayload,
                  currentEnvironment.requireComment
                )
              );
            });
          });
        };

        const hasDirtyField = (field) =>
          field && JSON.stringify(field).includes('true');

        if (saveFeatureType === SaveFeatureType.SCHEDULE) {
          const defaultValues = getDefaultValues(
            feature,
            currentEnvironment.requireComment
          );

          const prerequisitesList = [];
          if (hasDirtyField(dirtyFields.prerequisites)) {
            const defaultFeatureIds = new Map(
              defaultValues.prerequisites.map((o) => [
                o.featureId,
                o.variationId
              ])
            );
            const currentFeatureIds = new Map(
              data.prerequisites.map((v) => [v.featureId, v.variationId])
            );

            // Handle remove prerequisite
            defaultValues.prerequisites.forEach(({ featureId }) => {
              if (!currentFeatureIds.has(featureId)) {
                prerequisitesList.push(
                  createPrerequisiteChange(featureId, null, ChangeType.DELETE)
                );
              }
            });

            // Handle add & update prerequisite
            data.prerequisites.forEach(({ featureId, variationId }) => {
              if (!defaultFeatureIds.has(featureId)) {
                prerequisitesList.push(
                  createPrerequisiteChange(
                    featureId,
                    variationId,
                    ChangeType.CREATE
                  )
                );
              } else if (defaultFeatureIds.get(featureId) !== variationId) {
                prerequisitesList.push(
                  createPrerequisiteChange(
                    featureId,
                    variationId,
                    ChangeType.UPDATE
                  )
                );
              }
            });
          }

          const targets = [];
          if (hasDirtyField(dirtyFields.targets)) {
            defaultValues.targets.forEach((org, idx) => {
              const val = data.targets[idx];

              // Handle user removal from variation
              const removedUsers = org.users.filter(
                (u) => !val.users.includes(u)
              );
              if (removedUsers.length > 0) {
                targets.push(
                  createTargetChange(
                    org.variationId,
                    removedUsers,
                    ChangeType.DELETE
                  )
                );
              }

              // Handle user addition to variation
              const addedUsers = val.users.filter(
                (u) => !org.users.includes(u)
              );
              if (addedUsers.length > 0) {
                targets.push(
                  createTargetChange(
                    val.variationId,
                    addedUsers,
                    ChangeType.CREATE
                  )
                );
              }
            });
          }

          const rules = [];
          if (hasDirtyField(dirtyFields.rules)) {
            const orgRules = defaultValues.rules;
            const valRules = data.rules;

            const orgRuleIds = orgRules.map((r) => r.id);
            const valRuleIds = valRules.map((r) => r.id);

            const getRuleById = (rules, id) => rules.find((r) => r.id === id);

            const processRuleChange = (type, ruleData) => {
              const ruleChange = new RuleChange();
              ruleChange.setChangeType(type);
              ruleChange.setRule(ruleData);
              rules.push(ruleChange);
            };

            const processRemovedRules = () => {
              orgRules
                .filter((r) => !valRuleIds.includes(r.id))
                .forEach((r) => {
                  console.log('remove rule');
                  const rule = new Rule();
                  rule.setId(r.id);
                  processRuleChange(ChangeType.DELETE, rule);
                });
            };

            const processAddedRules = () => {
              valRules
                .filter((r) => !orgRuleIds.includes(r.id))
                .forEach((r) => {
                  console.log('add rule');
                  const rule = new Rule();
                  rule.setId(r.id);
                  rule.setStrategy(createStrategy(r.strategy));
                  rule.setClausesList(createClauses(r.clauses));
                  processRuleChange(ChangeType.CREATE, rule);
                });
            };

            const processUpdatedRules = () => {
              orgRuleIds
                .filter((id) => valRuleIds.includes(id))
                .forEach((rid) => {
                  let ruleUpdated = false;
                  const orgRule = getRuleById(orgRules, rid);
                  const valRule = getRuleById(valRules, rid);

                  const orgClauseIds = orgRule.clauses.map((c) => c.id);
                  const valClauseIds = valRule.clauses.map((c) => c.id);
                  const commonClauseIds = orgClauseIds.filter((id) =>
                    valClauseIds.includes(id)
                  );

                  if (orgClauseIds.some((id) => !valClauseIds.includes(id))) {
                    console.log('delete clause');
                    ruleUpdated = true;
                  }
                  if (valClauseIds.some((id) => !orgClauseIds.includes(id))) {
                    console.log('add clause');
                    ruleUpdated = true;
                  }

                  commonClauseIds.forEach((cid) => {
                    const orgClause = getRuleById(orgRule.clauses, cid);
                    const valClause = getRuleById(valRule.clauses, cid);

                    if (orgClause.attribute !== valClause.attribute) {
                      console.log('change attribute');
                      ruleUpdated = true;
                    }
                    if (orgClause.operator !== valClause.operator) {
                      console.log('change operator');
                      ruleUpdated = true;
                    }

                    const orgValues = new Set(orgClause.values);
                    const valValues = new Set(valClause.values);
                    const commonValues = [...orgValues].filter((v) =>
                      valValues.has(v)
                    );

                    [...orgValues]
                      .filter((v) => !commonValues.includes(v))
                      .forEach(() => {
                        console.log('remove clause value command');
                        ruleUpdated = true;
                      });
                    [...valValues]
                      .filter((v) => !commonValues.includes(v))
                      .forEach(() => {
                        console.log('add clause value command');
                        ruleUpdated = true;
                      });
                  });

                  if (
                    orgRule.strategy.option.value !==
                    valRule.strategy.option.value
                  ) {
                    console.log('change strategy command');
                    ruleUpdated = true;
                  }

                  if (
                    !deepEqual(
                      orgRule.strategy.rolloutStrategy,
                      valRule.strategy.rolloutStrategy
                    )
                  ) {
                    console.log('change rollout strategy command');
                    ruleUpdated = true;
                  }

                  if (
                    !deepEqual(
                      defaultValues.defaultStrategy,
                      data.defaultStrategy
                    )
                  ) {
                    console.log('change default strategy command');
                    ruleUpdated = true;
                  }

                  if (ruleUpdated) {
                    const rule = new Rule();
                    rule.setId(rid);
                    rule.setClausesList(createClauses(valRule.clauses));
                    rule.setStrategy(createStrategy(valRule.strategy));
                    processRuleChange(ChangeType.UPDATE, rule);
                  }
                });
            };

            processRemovedRules();
            processAddedRules();
            processUpdatedRules();
          }

          const updatePayload = {
            environmentId: currentEnvironment.id,
            id: featureId,
            comment: data.comment,
            enabled: dirtyFields.enabled ? data.enabled : undefined,
            prerequisitesList: prerequisitesList.length && prerequisitesList,
            targets: targets.length && targets,
            rules: rules.length && rules,
            defaultStrategy:
              hasDirtyField(dirtyFields.defaultStrategy) &&
              data.defaultStrategy,
            offVariation:
              hasDirtyField(dirtyFields.offVariation) && data.offVariation,
            resetSampling: data.resetSampling
          };

          await prepareUpdate(updateFeature, updatePayload);
        } else if (saveFeatureType === SaveFeatureType.UPDATE_NOW) {
          const commands: Array<Command> = [];
          const defaultValues = getDefaultValues(
            feature,
            currentEnvironment.requireComment
          );
          dirtyFields.enabled &&
            commands.push(...createEnabledCommands(defaultValues, data));
          dirtyFields.targets &&
            commands.push(
              ...createTargetCommands(defaultValues.targets, data.targets)
            );
          dirtyFields.rules &&
            commands.push(
              ...createRuleCommands(defaultValues.rules, data.rules)
            );
          dirtyFields.rules &&
            commands.push(
              ...createClauseCommands(defaultValues.rules, data.rules)
            );
          dirtyFields.rules &&
            commands.push(
              ...createStrategyCommands(defaultValues.rules, data.rules)
            );
          dirtyFields.defaultStrategy &&
            commands.push(
              ...createDefaultStrategyCommands(
                defaultValues.defaultStrategy,
                data.defaultStrategy
              )
            );
          dirtyFields.offVariation &&
            commands.push(
              ...createOffVariationCommands(
                defaultValues.offVariation,
                data.offVariation
              )
            );
          data.resetSampling && commands.push(createResetSampleSeedCommand());
          dirtyFields.prerequisites &&
            commands.push(
              ...createPrerequisitesCommands(
                defaultValues.prerequisites,
                data.prerequisites
              )
            );

          const updatePayload = {
            environmentId: currentEnvironment.id,
            id: featureId,
            comment: data.comment,
            commands: commands
          };
          await prepareUpdate(updateFeatureTargeting, updatePayload);
        }
      },
      [dispatch, dirtyFields, feature]
    );

    useEffect(() => {
      dispatch(
        listSegments({
          environmentId: currentEnvironment.id,
          cursor: ''
        })
      );
    }, [dispatch, currentEnvironment]);

    useEffect(() => {
      if (feature) {
        reset(getDefaultValues(feature, currentEnvironment.requireComment));
      }
    }, [feature]);

    const dirtyFieldsKeys = (obj) => {
      // Array containing the allowed keys
      const allowedKeys = ['enabled', 'comment', 'resetSampling'];

      // Check if any key in the object is not allowed
      for (const key in obj) {
        if (!allowedKeys.includes(key)) {
          return false; // Return false if any key is not allowed
        }
      }

      // Return true if all keys are allowed
      return true;
    };

    // Check if only switch is enabled/disabled or other fields are also changed
    // If only switch is enabled/disabled, then show Enable/Disable now and Schedule radio options in Confirm dialog
    // If other fields are also changed, then hide Enable/Disable now and Schedule radio options in Confirm dialog
    // const isSwitchEnabledConfirm = dirtyFieldsKeys(dirtyFields);

    return (
      <FormProvider {...methods}>
        <FeatureTargetingForm
          featureId={featureId}
          onOpenConfirmDialog={() => setIsConfirmDialogOpen(true)}
        />
        {isConfirmDialogOpen && (
          <FeatureConfirmDialog
            open={isConfirmDialogOpen}
            handleSubmit={(arg) => {
              handleSubmit((data) => handleUpdate(data, arg))();
            }}
            onClose={() => setIsConfirmDialogOpen(false)}
            title={f(messages.feature.confirm.title)}
            description={f(messages.feature.confirm.description)}
            displayResetSampling={true}
            featureId={featureId}
            // isSwitchEnabledConfirm={isSwitchEnabledConfirm}
            isEnabled={
              dirtyFields.enabled &&
              getDefaultValues(feature, currentEnvironment.requireComment)
                .enabled
            }
          />
        )}
      </FormProvider>
    );
  }
);

function createPrerequisiteChange(featureId, variationId, changeType) {
  const prerequisiteChange = new PrerequisiteChange();
  prerequisiteChange.setChangeType(changeType);

  const p = new Prerequisite();
  p.setFeatureId(featureId);
  if (variationId !== null) {
    p.setVariationId(variationId);
  }

  prerequisiteChange.setPrerequisite(p);
  return prerequisiteChange;
}

function createTargetChange(variationId, usersList, changeType) {
  const targetChange = new TargetChange();
  targetChange.setChangeType(changeType);

  const target = new Target();
  target.setVariation(variationId);
  target.setUsersList(usersList);

  targetChange.setTarget(target);
  return targetChange;
}

const createStrategyDefaultValue = (
  strategy: Strategy.AsObject,
  variations: Variation.AsObject[]
): StrategySchema => {
  return {
    option:
      strategy.type === Strategy.Type.FIXED
        ? {
            value: strategy.fixedStrategy.variation,
            label: createVariationLabel(
              variations.find((v) => v.id === strategy.fixedStrategy.variation)
            )
          }
        : {
            value: Strategy.Type.ROLLOUT.toString(),
            label: intl.formatMessage(
              messages.feature.strategy.selectRolloutPercentage
            )
          },
    rolloutStrategy: strategy.rolloutStrategy
      ? strategy.rolloutStrategy.variationsList.map((v) => {
          return {
            id: v.variation,
            percentage: v.weight / 1000
          };
        })
      : variations.map((v) => {
          return {
            id: v.id,
            percentage: 0
          };
        })
  };
};

const createClauseType = (
  cause: Clause.OperatorMap[keyof Clause.OperatorMap]
): RuleClauseType => {
  switch (cause) {
    case Clause.Operator.SEGMENT:
      return ruleClauseType.SEGMENT;
    case Clause.Operator.BEFORE:
    case Clause.Operator.AFTER:
      return ruleClauseType.DATE;
    case Clause.Operator.FEATURE_FLAG:
      return ruleClauseType.FEATURE_FLAG;
    default:
      return ruleClauseType.COMPARE;
  }
};

export const createEnabledCommands = (org, val): Command[] => {
  const commands: Command[] = [];
  if (org.enabled != val.enabled) {
    if (val.enabled) {
      const command = new EnableFeatureCommand();
      commands.push(
        createCommand({ message: command, name: 'EnableFeatureCommand' })
      );
    } else {
      const command = new DisableFeatureCommand();
      commands.push(
        createCommand({ message: command, name: 'DisableFeatureCommand' })
      );
    }
  }
  return commands;
};

export const createTargetCommands = (orgTargets, valTargets): Command[] => {
  const commands: Command[] = [];
  orgTargets.forEach((org, idx) => {
    const val = valTargets[idx];
    org.users
      .filter((u: string) => !val.users.includes(u))
      .forEach((u: string) => {
        const command = new RemoveUserFromVariationCommand();
        command.setId(org.variationId);
        command.setUser(u);
        commands.push(
          createCommand({
            message: command,
            name: 'RemoveUserFromVariationCommand'
          })
        );
      });
    val.users
      .filter((u: string) => !org.users.includes(u))
      .forEach((u: string) => {
        const command = new AddUserToVariationCommand();
        command.setId(org.variationId);
        command.setUser(u);
        commands.push(
          createCommand({
            message: command,
            name: 'AddUserToVariationCommand'
          })
        );
      });
  });
  return commands;
};

export function createRuleCommands(org, val): Command[] {
  const commands: Array<Command> = [];
  const orgIds = org.map((r) => r.id);
  const valIds = val.map((r) => r.id);

  org
    .filter((r) => !valIds.includes(r.id))
    .forEach((r) => {
      console.log('delete rule command');
      const command = new DeleteRuleCommand();
      command.setId(r.id);
      commands.push(
        createCommand({ message: command, name: 'DeleteRuleCommand' })
      );
    });

  val
    .filter((r) => !orgIds.includes(r.id))
    .forEach((r) => {
      console.log('add rule command');
      const command = new AddRuleCommand();
      command.setRule(createRule(r));
      commands.push(
        createCommand({ message: command, name: 'AddRuleCommand' })
      );
    });

  let orderChanged = false;
  const orgIdsAfterDeletedIdsRemoved = orgIds.filter((orgId) =>
    valIds.includes(orgId)
  );

  // check if any rule is deleted
  if (org.find((r) => !valIds.includes(r.id))) {
    // check if order changed
    if (
      !orderChanged &&
      valIds.slice(0, orgIdsAfterDeletedIdsRemoved.length).toString() !==
        orgIdsAfterDeletedIdsRemoved.toString()
    ) {
      console.log('rule deleted order changed');
      orderChanged = true;
      commands.push(createChangeRulesOrderCommand(valIds));
    }
  }

  // check if any rule is added
  if (val.find((r) => !orgIds.includes(r.id))) {
    // check if order changed
    if (
      !orderChanged &&
      orgIdsAfterDeletedIdsRemoved.toString() !==
        valIds.slice(0, orgIdsAfterDeletedIdsRemoved.length).toString()
    ) {
      console.log('rule added order changed');
      orderChanged = true;
      commands.push(createChangeRulesOrderCommand(valIds));
    }
  }

  // check if only rule order changed
  if (
    !orderChanged &&
    orgIds.length === valIds.length &&
    orgIds.every((orgId) => valIds.includes(orgId)) &&
    orgIds.toString() !== valIds.toString()
  ) {
    console.log('rule order changed');
    commands.push(createChangeRulesOrderCommand(valIds));
  }
  return commands;
}

const createChangeRulesOrderCommand = (valIds: string[]): Command => {
  const command = new ChangeRulesOrderCommand();
  command.setRuleIdsList(valIds);
  return createCommand({
    message: command,
    name: 'ChangeRulesOrderCommand'
  });
};

const createRule = (rule): Rule => {
  const r = new Rule();
  r.setId(rule.id);
  r.setStrategy(createStrategy(rule.strategy));
  r.setClausesList(createClauses(rule.clauses));
  return r;
};

const createStrategy = (strategy: StrategySchema): Strategy => {
  const s = new Strategy();
  if (strategy.option.value == Strategy.Type.ROLLOUT.toString()) {
    const variations = strategy.rolloutStrategy.map((rs) => {
      const v = new RolloutStrategy.Variation();
      v.setVariation(rs.id);
      v.setWeight(rs.percentage * 1000);
      return v;
    });
    const rs = new RolloutStrategy();
    rs.setVariationsList(variations);
    s.setType(Strategy.Type.ROLLOUT);
    s.setRolloutStrategy(rs);
    return s;
  }
  s.setType(Strategy.Type.FIXED);
  const fs = new FixedStrategy();
  fs.setVariation(strategy.option.value as string);
  s.setFixedStrategy(fs);
  return s;
};

const createClauses = (clauses: RuleClauseSchema[]): Clause[] => {
  const cs: Clause[] = [];
  clauses.forEach((c) => {
    cs.push(createClause(c));
  });
  return cs;
};

const createClause = (clause: RuleClauseSchema): Clause => {
  const c = new Clause();
  c.setId(clause.id);
  c.setAttribute(clause.attribute);
  c.setOperator(
    Number(clause.operator) as Clause.OperatorMap[keyof Clause.OperatorMap]
  );
  c.setValuesList(clause.values);
  return c;
};

export const createClauseCommands = (
  orgRules: RuleSchema[],
  valRules: RuleSchema[]
): Array<Command> => {
  const commands: Command[] = [];
  const orgRuleIds = orgRules.filter((r) => r.id).map((r) => r.id);
  const valRuleIds = valRules.filter((r) => r.id).map((r) => r.id);
  // Intersection of org and val rules.
  const rulesIds = orgRuleIds.filter((id) => valRuleIds.includes(id));
  rulesIds.forEach((rid) => {
    const orgRule = orgRules.find((r) => r.id === rid);
    const valRule = valRules.find((r) => r.id === rid);
    const orgClauseIds = orgRule.clauses.filter((c) => c.id).map((c) => c.id);
    const valClauseIds = valRule.clauses.filter((c) => c.id).map((c) => c.id);
    // Intersection of org and val clauses.
    const clauseIds = orgClauseIds.filter((id) => valClauseIds.includes(id));
    orgRule.clauses
      .filter((c) => !clauseIds.includes(c.id))
      .forEach((c) => {
        console.log('delete clause command', c);
        const command = new DeleteClauseCommand();
        command.setRuleId(rid);
        command.setId(c.id);
        commands.push(
          createCommand({ message: command, name: 'DeleteClauseCommand' })
        );
      });
    valRule.clauses
      .filter((c) => !clauseIds.includes(c.id))
      .forEach((c) => {
        console.log('add clause command', c);
        const command = new AddClauseCommand();
        command.setRuleId(rid);
        command.setClause(createClause(c));
        commands.push(
          createCommand({ message: command, name: 'AddClauseCommand' })
        );
      });
    clauseIds.forEach((cid) => {
      const orgClause = orgRule.clauses.find((c) => c.id === cid);
      const valClause = valRule.clauses.find((c) => c.id === cid);
      commands.push(
        ...createClauseAttributeCommands(rid, orgClause, valClause)
      );
      commands.push(...createClauseValueCommands(rid, orgClause, valClause));
      commands.push(...createClauseOperatorCommands(rid, orgClause, valClause));
    });
  });
  return commands;
};

function createClauseAttributeCommands(
  ruleId: string,
  orgClause,
  valClause
): Command[] {
  const commands: Command[] = [];
  if (orgClause.attribute !== valClause.attribute) {
    const command = new ChangeClauseAttributeCommand();
    command.setRuleId(ruleId);
    command.setId(orgClause.id);
    command.setAttribute(valClause.attribute);
    commands.push(
      createCommand({ message: command, name: 'ChangeClauseAttributeCommand' })
    );
  }
  return commands;
}

function createClauseValueCommands(
  ruleId: string,
  orgClause: RuleClauseSchema,
  valClause: RuleClauseSchema
): Command[] {
  const commands: Command[] = [];
  // Intersection of org and val values.
  const orgValues = orgClause.values;
  const valValues = valClause.values;
  const values = orgValues.filter((v) => valValues.includes(v));
  orgValues
    .filter((v) => !values.includes(v))
    .forEach((v) => {
      const command = new RemoveClauseValueCommand();
      command.setId(orgClause.id);
      command.setRuleId(ruleId);
      command.setValue(String(v));
      commands.push(
        createCommand({ message: command, name: 'RemoveClauseValueCommand' })
      );
    });
  valValues
    .filter((v) => !values.includes(v))
    .forEach((v) => {
      const command = new AddClauseValueCommand();
      command.setRuleId(ruleId);
      command.setId(orgClause.id);
      command.setValue(String(v));
      commands.push(
        createCommand({ message: command, name: 'AddClauseValueCommand' })
      );
    });
  return commands;
}

const createClauseOperatorCommands = (
  ruleId: string,
  orgClause: RuleClauseSchema,
  valClause: RuleClauseSchema
): Command[] => {
  const commands: Command[] = [];
  if (orgClause.operator != valClause.operator) {
    const command = new ChangeClauseOperatorCommand();
    command.setRuleId(ruleId);
    command.setId(orgClause.id);
    command.setOperator(
      Number(valClause.operator) as Clause.OperatorMap[keyof Clause.OperatorMap]
    );
    commands.push(
      createCommand({ message: command, name: 'ChangeClauseOperatorCommand' })
    );
  }
  return commands;
};

const createStrategyCommands = (
  orgRules: RuleSchema[],
  valRules: RuleSchema[]
): Command[] => {
  const commands: Command[] = [];
  const orgRuleIds = orgRules.filter((r) => r.id).map((r) => r.id);
  const valRuleIds = valRules.filter((r) => r.id).map((r) => r.id);
  // Intersection of org and val rules.
  const ruleIds = orgRuleIds.filter((id) => valRuleIds.includes(id));
  ruleIds.forEach((rid) => {
    const orgRule = orgRules.find((r) => r.id === rid);
    const valRule = valRules.find((r) => r.id === rid);
    if (orgRule.strategy.option.value != valRule.strategy.option.value) {
      if (
        orgRule.strategy.option.value == Strategy.Type.ROLLOUT.toString() ||
        valRule.strategy.option.value == Strategy.Type.ROLLOUT.toString()
      ) {
        console.log('change rule strategy command');
        const command = new ChangeRuleStrategyCommand();
        command.setRuleId(rid);
        command.setStrategy(createStrategy(valRule.strategy));
        commands.push(
          createCommand({
            message: command,
            name: 'ChangeRuleStrategyCommand'
          })
        );
        return;
      }
      const command = new ChangeFixedStrategyCommand();
      command.setRuleId(rid);
      command.setStrategy(createStrategy(valRule.strategy).getFixedStrategy());
      commands.push(
        createCommand({
          message: command,
          name: 'ChangeFixedStrategyCommand'
        })
      );
      return;
    }
    if (
      !deepEqual(
        orgRule.strategy.rolloutStrategy,
        valRule.strategy.rolloutStrategy
      )
    ) {
      const command = new ChangeRolloutStrategyCommand();
      command.setRuleId(rid);
      command.setStrategy(
        createStrategy(valRule.strategy).getRolloutStrategy()
      );
      commands.push(
        createCommand({
          message: command,
          name: 'ChangeRolloutStrategyCommand'
        })
      );
    }
  });
  return commands;
};

export function createDefaultStrategyCommands(org, val): Command[] {
  const commands: Command[] = [];
  if (!deepEqual(org, val)) {
    const command = new ChangeDefaultStrategyCommand();
    command.setStrategy(createStrategy(val));
    commands.push(
      createCommand({ message: command, name: 'ChangeDefaultStrategyCommand' })
    );
  }
  return commands;
}

export function createOffVariationCommands(org, val): Command[] {
  const commands: Command[] = [];
  if (org.value !== val.value) {
    const command = new ChangeOffVariationCommand();
    command.setId(val.value);
    commands.push(
      createCommand({ message: command, name: 'ChangeOffVariationCommand' })
    );
  }
  return commands;
}

export function createPrerequisitesCommands(org, val): Command[] {
  const commands: Array<Command> = [];

  // handle remove feature
  org.filter((o) => {
    if (!val.some((v) => v.featureId === o.featureId)) {
      const command = new RemovePrerequisiteCommand();
      command.setFeatureId(o.featureId);
      commands.push(
        createCommand({ message: command, name: 'RemovePrerequisiteCommand' })
      );
    }
  });

  // handle add feature
  val.filter((v) => {
    if (!org.some((o) => o.featureId === v.featureId)) {
      const command = new AddPrerequisiteCommand();
      command.setPrerequisite(createPrerequisite(v));
      commands.push(
        createCommand({ message: command, name: 'AddPrerequisiteCommand' })
      );
    }
  });

  // handle update variation
  val.forEach((v) => {
    if (
      org.some(
        (o) => o.featureId === v.featureId && o.variationId !== v.variationId
      )
    ) {
      const command = new ChangePrerequisiteVariationCommand();
      command.setPrerequisite(createPrerequisite(v));
      commands.push(
        createCommand({
          message: command,
          name: 'ChangePrerequisiteVariationCommand'
        })
      );
    }
  });

  return commands;
}

const createPrerequisite = (prerequisite): Prerequisite => {
  const p = new Prerequisite();
  p.setFeatureId(prerequisite.featureId);
  p.setVariationId(prerequisite.variationId);
  return p;
};

export function createResetSampleSeedCommand(): Command {
  const command = new ResetSamplingSeedCommand();
  return createCommand({ message: command, name: 'ResetSamplingSeedCommand' });
}
