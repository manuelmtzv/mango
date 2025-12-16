export type PaginateConfig = {
  count: number;
  page?: number;
  limit?: number;
};

export function paginate(config: PaginateConfig) {
  const { page = 1, limit = 10 } = config;

  const offset = (page - 1) * limit;

  return {
    config: () => ({
      take: limit,
      skip: offset,
    }),

    maxPage: Math.ceil(config.count / limit),
  };
}
