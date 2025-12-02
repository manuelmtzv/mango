export interface IPaginatedReponse<T> {
  data: T[];
  meta: {
    page: number;
    limit: number;
    count: number;
    totalPages: number;
  };
}
