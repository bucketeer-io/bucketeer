import * as yup from 'yup';

export const formSchema = yup.object().shape({
  prerequisitesRules: yup
    .array()
    .required()
    .of(
      yup.object().shape({
        index: yup.number().required(),
        rules: yup
          .array()
          .required()
          .of(
            yup.object().shape({
              featureFlag: yup.string().required('This field is required.'),
              variation: yup.string().required('This field is required.')
            })
          )
      })
    ),
  targetIndividualRules: yup
    .array()
    .required()
    .of(
      yup.object().shape({
        on: yup.array().required('This field is required.'),
        off: yup.array().required('This field is required.')
      })
    ),
  targetSegmentRules: yup
    .array()
    .required()
    .of(
      yup.object().shape({
        index: yup.number().required(),
        rules: yup
          .array()
          .required()
          .of(
            yup.object().shape({
              conditions: yup
                .array()
                .required()
                .of(
                  yup.object().shape({
                    situation: yup
                      .string()
                      .oneOf([
                        'compare',
                        'user-segment',
                        'date',
                        'feature-flag'
                      ])
                      .required('This field is required.'),
                    conditioner: yup
                      .string()
                      .required('This field is required.'),
                    firstValue: yup
                      .string()
                      .test('required', (value, context) => {
                        const situation =
                          context.from && context.from[0].value.situation;
                        if (!value && situation === 'compare')
                          return context.createError({
                            message: `This field is required.`,
                            path: context.path
                          });

                        return true;
                      }),
                    secondValue: yup
                      .string()
                      .test('required', (value, context) => {
                        const situation =
                          context.from && context.from[0].value.situation;
                        if (!value && situation === 'compare')
                          return context.createError({
                            message: `This field is required.`,
                            path: context.path
                          });

                        return true;
                      }),
                    value: yup.string().test('required', (value, context) => {
                      const situation =
                        context.from && context.from[0].value.situation;
                      if (
                        !value &&
                        ['user-segment', 'date'].includes(situation)
                      )
                        return context.createError({
                          message: `This field is required.`,
                          path: context.path
                        });

                      return true;
                    }),
                    date: yup.string().test('required', (value, context) => {
                      const situation =
                        context.from && context.from[0].value.situation;
                      if (!value && situation === 'date')
                        return context.createError({
                          message: `This field is required.`,
                          path: context.path
                        });
                      return true;
                    }),
                    flagId: yup.string().test('required', (value, context) => {
                      const situation =
                        context.from && context.from[0].value.situation;
                      if (!value && situation === 'feature-flag')
                        return context.createError({
                          message: `This field is required.`,
                          path: context.path
                        });

                      return true;
                    }),
                    variation: yup
                      .string()
                      .test('required', (value, context) => {
                        const situation =
                          context.from && context.from[0].value.situation;
                        if (!value && situation === 'feature-flag')
                          return context.createError({
                            message: `This field is required.`,
                            path: context.path
                          });

                        return true;
                      })
                  })
                ),
              variation: yup.boolean().required('This field is required.')
            })
          )
      })
    )
});
