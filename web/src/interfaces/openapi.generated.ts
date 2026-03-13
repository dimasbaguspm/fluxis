/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface DomainAuthLoginModel {
  /** @example "user@example.com" */
  email: string;
  /** @example "s3cr3tP@ssword" */
  password: string;
}

export interface DomainAuthModel {
  /** @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." */
  accessToken?: string;
  /** @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." */
  refreshToken?: string;
}

export interface DomainAuthRefreshModel {
  /** @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." */
  accessToken: string;
  /** @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." */
  refreshToken: string;
}

export interface DomainAuthRegisterModel {
  /**
   * @minLength 1
   * @example "John Doe"
   */
  displayName?: string;
  /** @example "user@example.com" */
  email: string;
  /** @example "s3cr3tP@ssword" */
  password: string;
}

export interface DomainBoardColumnCreateModel {
  /** @minLength 1 */
  name: string;
}

export interface DomainBoardColumnModel {
  boardId?: string;
  createdAt?: string;
  id?: string;
  /** @minLength 1 */
  name: string;
  position?: number;
  updatedAt?: string;
}

export interface DomainBoardColumnUpdateModel {
  /** @minLength 1 */
  name?: string;
}

export interface DomainBoardColumnsPagedModel {
  items?: DomainBoardColumnModel[];
  pageNumber?: number;
  pageSize?: number;
  totalCount?: number;
  totalPages?: number;
}

export interface DomainBoardCreateModel {
  /** @minLength 1 */
  name: string;
  sprintId: string;
}

export interface DomainBoardModel {
  createdAt?: string;
  id?: string;
  /** @minLength 1 */
  name: string;
  position?: number;
  sprintId?: string;
  updatedAt?: string;
}

export interface DomainBoardUpdateModel {
  /** @minLength 1 */
  name?: string;
  sprintId?: string;
}

export interface DomainBoardsPagedModel {
  items?: DomainBoardModel[];
  pageNumber?: number;
  pageSize?: number;
  totalCount?: number;
  totalPages?: number;
}

export interface DomainOrganisationCreateModel {
  /** @minLength 1 */
  name: string;
}

export interface DomainOrganisationMemberCreateModel {
  role: "admin" | "member" | "viewer";
  userId: string;
}

export interface DomainOrganisationMemberModel {
  email?: string;
  joinedAt?: string;
  name?: string;
  role?: string;
  userId?: string;
}

export interface DomainOrganisationMemberUpdateModel {
  role: "admin" | "member" | "viewer";
}

export interface DomainOrganisationMembersPagedModel {
  items?: DomainOrganisationMemberModel[];
  pageNumber?: number;
  pageSize?: number;
  totalCount?: number;
  totalPages?: number;
}

export interface DomainOrganisationModel {
  createdAt?: string;
  id: string;
  /** @minLength 1 */
  name?: string;
  slug?: string;
  totalMembers?: number;
  updatedAt?: string;
}

export interface DomainOrganisationPagedModel {
  items?: DomainOrganisationModel[];
  page?: number;
  pageSize?: number;
  totalCount?: number;
  totalPages?: number;
}

export interface DomainOrganisationUpdateModel {
  /** @minLength 1 */
  name?: string;
}

export interface DomainProjectCreateModel {
  description?: string;
  /**
   * @minLength 1
   * @maxLength 10
   */
  key: string;
  /**
   * @minLength 1
   * @maxLength 100
   */
  name: string;
  visibility: "public" | "private";
}

export interface DomainProjectModel {
  createdAt?: string;
  description?: string;
  id: string;
  /** @minLength 1 */
  key: string;
  /** @minLength 1 */
  name: string;
  orgId: string;
  updatedAt?: string;
  visibility: "public" | "private";
}

export interface DomainProjectUpdateModel {
  description?: string;
  /**
   * @minLength 1
   * @maxLength 100
   */
  name?: string;
}

export interface DomainProjectVisibilityModel {
  visibility: "public" | "private";
}

export interface DomainProjectsPagedModel {
  items?: DomainProjectModel[];
  pageNumber?: number;
  pageSize?: number;
  totalCount?: number;
  totalPages?: number;
}

export interface DomainSprintCreateModel {
  goal?: string;
  /** @minLength 1 */
  name: string;
  plannedCompletedAt?: string;
  plannedStartedAt?: string;
  projectId: string;
  status?: "planned" | "active" | "completed";
}

export interface DomainSprintModel {
  completedAt?: string;
  createdAt?: string;
  goal?: string;
  id?: string;
  /** @minLength 1 */
  name: string;
  plannedCompletedAt?: string;
  plannedStartedAt?: string;
  projectId?: string;
  startedAt?: string;
  status: "planned" | "active" | "completed";
  updatedAt?: string;
}

export interface DomainSprintUpdateModel {
  goal?: string;
  /** @minLength 1 */
  name?: string;
  plannedCompletedAt?: string;
  plannedStartedAt?: string;
  status?: "planned" | "active" | "completed";
}

export interface DomainSprintsPagedModel {
  items?: DomainSprintModel[];
  pageNumber?: number;
  pageSize?: number;
  totalCount?: number;
  totalPages?: number;
}

export interface DomainTicketBoardMoveModel {
  boardColumnId: string;
  boardId: string;
}

export interface DomainTicketCreateModel {
  assigneeId?: string;
  description?: string;
  dueDate?: string;
  priority: "low" | "medium" | "high" | "critical";
  sprintId?: string;
  /** @min 0 */
  storyPoints?: number;
  /**
   * @minLength 1
   * @maxLength 255
   */
  title: string;
  type: "bug" | "story" | "task" | "epic";
}

export interface DomainTicketModel {
  assigneeId?: string;
  boardColumnId?: string;
  boardId?: string;
  createdAt?: string;
  description?: string;
  dueDate?: string;
  epicId?: string;
  id: string;
  key?: string;
  parentId?: string;
  priority?: string;
  projectId: string;
  reporterId?: string;
  sprintId?: string;
  storyPoints?: number;
  ticketNumber?: number;
  title?: string;
  type?: string;
  updatedAt?: string;
}

export interface DomainTicketUpdateModel {
  assigneeId?: string;
  description?: string;
  dueDate?: string;
  priority?: "low" | "medium" | "high" | "critical";
  sprintId?: string;
  /** @min 0 */
  storyPoints?: number;
  /**
   * @minLength 1
   * @maxLength 255
   */
  title?: string;
  type?: "bug" | "story" | "task" | "epic";
}

export interface DomainTicketsPagedModel {
  items?: DomainTicketModel[];
  pageNumber?: number;
  pageSize?: number;
  totalCount?: number;
  totalPages?: number;
}

export interface DomainUserModel {
  createdAt?: string;
  /** @example "John Doe" */
  displayName?: string;
  /** @example "user@example.com" */
  email?: string;
  /** @example "550e8400-e29b-41d4-a716-446655440000" */
  id: string;
  password?: string;
  updatedAt?: string;
}

export interface HttpxErrBlock {
  /** machine-readable e.g. "email_taken" */
  code?: string;
  message?: string;
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseFormat;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<
  FullRequestParams,
  "body" | "method" | "query" | "path"
>;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<RequestParams | void> | RequestParams | void;
  customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown>
  extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  JsonApi = "application/vnd.api+json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "http://localhost:8080";
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private abortControllers = new Map<CancelToken, AbortController>();
  private customFetch = (...fetchParams: Parameters<typeof fetch>) =>
    fetch(...fetchParams);

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key);
    return `${encodedKey}=${encodeURIComponent(typeof value === "number" ? value : `${value}`)}`;
  }

  protected addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key]);
  }

  protected addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];
    return value.map((v: any) => this.encodeQueryParam(key, v)).join("&");
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter(
      (key) => "undefined" !== typeof query[key],
    );
    return keys
      .map((key) =>
        Array.isArray(query[key])
          ? this.addArrayQueryParam(query, key)
          : this.addQueryParam(query, key),
      )
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.JsonApi]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.Text]: (input: any) =>
      input !== null && typeof input !== "string"
        ? JSON.stringify(input)
        : input,
    [ContentType.FormData]: (input: any) => {
      if (input instanceof FormData) {
        return input;
      }

      return Object.keys(input || {}).reduce((formData, key) => {
        const property = input[key];
        formData.append(
          key,
          property instanceof Blob
            ? property
            : typeof property === "object" && property !== null
              ? JSON.stringify(property)
              : `${property}`,
        );
        return formData;
      }, new FormData());
    },
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  protected mergeRequestParams(
    params1: RequestParams,
    params2?: RequestParams,
  ): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected createAbortSignal = (
    cancelToken: CancelToken,
  ): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = async <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format,
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.baseApiParams.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];
    const responseFormat = format || requestParams.format;

    return this.customFetch(
      `${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`,
      {
        ...requestParams,
        headers: {
          ...(requestParams.headers || {}),
          ...(type && type !== ContentType.FormData
            ? { "Content-Type": type }
            : {}),
        },
        signal:
          (cancelToken
            ? this.createAbortSignal(cancelToken)
            : requestParams.signal) || null,
        body:
          typeof body === "undefined" || body === null
            ? null
            : payloadFormatter(body),
      },
    ).then(async (response) => {
      const r = response as HttpResponse<T, E>;
      r.data = null as unknown as T;
      r.error = null as unknown as E;

      const responseToParse = responseFormat ? response.clone() : response;
      const data = !responseFormat
        ? r
        : await responseToParse[responseFormat]()
            .then((data) => {
              if (r.ok) {
                r.data = data;
              } else {
                r.error = data;
              }
              return r;
            })
            .catch((e) => {
              r.error = e;
              return r;
            });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title Fluxis API
 * @version 1.0
 * @license MIT
 * @baseUrl http://localhost:8080
 * @contact Fluxis Support (https://github.com/dimasbaguspm/fluxis)
 *
 * Personal finance management API
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  auth = {
    /**
     * @description Authenticates a user and returns access/refresh tokens
     *
     * @tags auth
     * @name LoginCreate
     * @summary Login with email and password
     * @request POST:/auth/login
     */
    loginCreate: (body: DomainAuthLoginModel, params: RequestParams = {}) =>
      this.request<DomainAuthModel, HttpxErrBlock>({
        path: `/auth/login`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Issues a new access token using a valid refresh token
     *
     * @tags auth
     * @name RefreshCreate
     * @summary Rotate access token
     * @request POST:/auth/refresh
     */
    refreshCreate: (body: DomainAuthRefreshModel, params: RequestParams = {}) =>
      this.request<DomainAuthModel, HttpxErrBlock>({
        path: `/auth/refresh`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Creates a new user account and returns access/refresh tokens
     *
     * @tags auth
     * @name RegisterCreate
     * @summary Register a new user
     * @request POST:/auth/register
     */
    registerCreate: (
      body: DomainAuthRegisterModel,
      params: RequestParams = {},
    ) =>
      this.request<DomainAuthModel, HttpxErrBlock>({
        path: `/auth/register`,
        method: "POST",
        body: body,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  boards = {
    /**
     * @description Returns all boards in a sprint with pagination
     *
     * @tags board
     * @name BoardsList
     * @summary List boards
     * @request GET:/boards
     * @secure
     */
    boardsList: (
      query?: {
        id?: string[];
        name?: string;
        /** @min 1 */
        pageNumber?: number;
        /**
         * @min 1
         * @max 100
         */
        pageSize?: number;
        sprintId?: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<DomainBoardsPagedModel, HttpxErrBlock>({
        path: `/boards`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Creates a new board in a sprint
     *
     * @tags board
     * @name BoardsCreate
     * @summary Create a board
     * @request POST:/boards
     * @secure
     */
    boardsCreate: (body: DomainBoardCreateModel, params: RequestParams = {}) =>
      this.request<DomainBoardModel, HttpxErrBlock>({
        path: `/boards`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Reorder boards within a sprint (positions determined by array order)
     *
     * @tags board
     * @name ReorderPartialUpdate
     * @summary Reorder boards
     * @request PATCH:/boards/reorder
     * @secure
     */
    reorderPartialUpdate: (
      query: {
        /** Sprint ID */
        sprintId: string;
      },
      body: string[],
      params: RequestParams = {},
    ) =>
      this.request<DomainBoardModel[], HttpxErrBlock>({
        path: `/boards/reorder`,
        method: "PATCH",
        query: query,
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Returns a single board by ID
     *
     * @tags board
     * @name BoardsDetail
     * @summary Get a board
     * @request GET:/boards/{boardId}
     * @secure
     */
    boardsDetail: (boardId: string, params: RequestParams = {}) =>
      this.request<DomainBoardModel, HttpxErrBlock>({
        path: `/boards/${boardId}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Deletes a board
     *
     * @tags board
     * @name BoardsDelete
     * @summary Delete a board
     * @request DELETE:/boards/{boardId}
     * @secure
     */
    boardsDelete: (boardId: string, params: RequestParams = {}) =>
      this.request<void, HttpxErrBlock>({
        path: `/boards/${boardId}`,
        method: "DELETE",
        secure: true,
        ...params,
      }),

    /**
     * @description Updates board details
     *
     * @tags board
     * @name BoardsPartialUpdate
     * @summary Update a board
     * @request PATCH:/boards/{boardId}
     * @secure
     */
    boardsPartialUpdate: (
      boardId: string,
      body: DomainBoardUpdateModel,
      params: RequestParams = {},
    ) =>
      this.request<DomainBoardModel, HttpxErrBlock>({
        path: `/boards/${boardId}`,
        method: "PATCH",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Returns all columns in a board with pagination
     *
     * @tags board
     * @name ColumnsList
     * @summary List board columns
     * @request GET:/boards/{boardId}/columns
     * @secure
     */
    columnsList: (
      boardId: string,
      query?: {
        boardId?: string[];
        id?: string[];
        name?: string;
        /** @min 1 */
        pageNumber?: number;
        /**
         * @min 1
         * @max 100
         */
        pageSize?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<DomainBoardColumnsPagedModel, HttpxErrBlock>({
        path: `/boards/${boardId}/columns`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Creates a new column in a board (position is auto-calculated)
     *
     * @tags board
     * @name ColumnsCreate
     * @summary Create a board column
     * @request POST:/boards/{boardId}/columns
     * @secure
     */
    columnsCreate: (
      boardId: string,
      body: DomainBoardColumnCreateModel,
      params: RequestParams = {},
    ) =>
      this.request<DomainBoardColumnModel, HttpxErrBlock>({
        path: `/boards/${boardId}/columns`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Reorder columns within a board
     *
     * @tags board
     * @name ColumnsReorderPartialUpdate
     * @summary Reorder board columns
     * @request PATCH:/boards/{boardId}/columns/reorder
     * @secure
     */
    columnsReorderPartialUpdate: (
      boardId: string,
      body: string[],
      params: RequestParams = {},
    ) =>
      this.request<DomainBoardColumnModel[], HttpxErrBlock>({
        path: `/boards/${boardId}/columns/reorder`,
        method: "PATCH",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Deletes a column from a board
     *
     * @tags board
     * @name ColumnsDelete
     * @summary Delete a board column
     * @request DELETE:/boards/{boardId}/columns/{boardColumnId}
     * @secure
     */
    columnsDelete: (
      boardId: string,
      boardColumnId: string,
      params: RequestParams = {},
    ) =>
      this.request<void, HttpxErrBlock>({
        path: `/boards/${boardId}/columns/${boardColumnId}`,
        method: "DELETE",
        secure: true,
        ...params,
      }),

    /**
     * @description Updates column name (use reorder endpoint for position changes)
     *
     * @tags board
     * @name ColumnsPartialUpdate
     * @summary Update a board column
     * @request PATCH:/boards/{boardId}/columns/{boardColumnId}
     * @secure
     */
    columnsPartialUpdate: (
      boardId: string,
      boardColumnId: string,
      body: DomainBoardColumnUpdateModel,
      params: RequestParams = {},
    ) =>
      this.request<DomainBoardColumnModel, HttpxErrBlock>({
        path: `/boards/${boardId}/columns/${boardColumnId}`,
        method: "PATCH",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  orgs = {
    /**
     * @description Returns paginated organisations with optional filtering and sorting
     *
     * @tags org
     * @name OrgsList
     * @summary List organisations with pagination
     * @request GET:/orgs
     * @secure
     */
    orgsList: (
      query?: {
        id?: string[];
        name?: string[];
        /** @min 1 */
        pageNumber?: number;
        /**
         * @min 1
         * @max 100
         */
        pageSize?: number;
        sortBy?: "name" | "createdAt" | "updatedAt";
        sortOrder?: "asc" | "desc";
      },
      params: RequestParams = {},
    ) =>
      this.request<DomainOrganisationPagedModel, HttpxErrBlock>({
        path: `/orgs`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Creates a new organisation
     *
     * @tags org
     * @name OrgsCreate
     * @summary Create an organisation
     * @request POST:/orgs
     * @secure
     */
    orgsCreate: (
      body: DomainOrganisationCreateModel,
      params: RequestParams = {},
    ) =>
      this.request<DomainOrganisationModel, HttpxErrBlock>({
        path: `/orgs`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Returns a single organisation by ID
     *
     * @tags org
     * @name OrgsDetail
     * @summary Get an organisation
     * @request GET:/orgs/{id}
     * @secure
     */
    orgsDetail: (id: string, params: RequestParams = {}) =>
      this.request<DomainOrganisationModel, HttpxErrBlock>({
        path: `/orgs/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Soft-deletes an organisation by ID
     *
     * @tags org
     * @name OrgsDelete
     * @summary Delete an organisation
     * @request DELETE:/orgs/{id}
     * @secure
     */
    orgsDelete: (id: string, params: RequestParams = {}) =>
      this.request<void, HttpxErrBlock>({
        path: `/orgs/${id}`,
        method: "DELETE",
        secure: true,
        ...params,
      }),

    /**
     * @description Updates an organisation's name
     *
     * @tags org
     * @name OrgsPartialUpdate
     * @summary Update an organisation
     * @request PATCH:/orgs/{id}
     * @secure
     */
    orgsPartialUpdate: (
      id: string,
      body: DomainOrganisationUpdateModel,
      params: RequestParams = {},
    ) =>
      this.request<DomainOrganisationModel, HttpxErrBlock>({
        path: `/orgs/${id}`,
        method: "PATCH",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Returns paginated members of an organisation with optional filtering
     *
     * @tags org
     * @name MembersList
     * @summary List organisation members with pagination
     * @request GET:/orgs/{id}/members
     * @secure
     */
    membersList: (
      id: string,
      query?: {
        displayName?: string;
        email?: string;
        /** @min 1 */
        pageNumber?: number;
        /**
         * @min 1
         * @max 100
         */
        pageSize?: number;
        userId?: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<DomainOrganisationMembersPagedModel, HttpxErrBlock>({
        path: `/orgs/${id}/members`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Adds a user to an organisation with a given role
     *
     * @tags org
     * @name MembersCreate
     * @summary Add a member to an organisation
     * @request POST:/orgs/{id}/members
     * @secure
     */
    membersCreate: (
      id: string,
      body: DomainOrganisationMemberCreateModel,
      params: RequestParams = {},
    ) =>
      this.request<void, HttpxErrBlock>({
        path: `/orgs/${id}/members`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description Delete a user from an organisation
     *
     * @tags org
     * @name MembersDelete
     * @summary Delete a member from an organsiation
     * @request DELETE:/orgs/{id}/members/{userId}
     * @secure
     */
    membersDelete: (id: string, userId: string, params: RequestParams = {}) =>
      this.request<void, HttpxErrBlock>({
        path: `/orgs/${id}/members/${userId}`,
        method: "DELETE",
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description Updates the role of a member within an organisation
     *
     * @tags org
     * @name MembersPartialUpdate
     * @summary Update a member's role
     * @request PATCH:/orgs/{id}/members/{userId}
     * @secure
     */
    membersPartialUpdate: (
      id: string,
      userId: string,
      body: DomainOrganisationMemberUpdateModel,
      params: RequestParams = {},
    ) =>
      this.request<void, HttpxErrBlock>({
        path: `/orgs/${id}/members/${userId}`,
        method: "PATCH",
        body: body,
        secure: true,
        type: ContentType.Json,
        ...params,
      }),
  };
  projects = {
    /**
     * @description Returns paginated projects in an organisation with optional filtering
     *
     * @tags project
     * @name ProjectsList
     * @summary List projects with pagination
     * @request GET:/projects
     * @secure
     */
    projectsList: (
      query?: {
        id?: string[];
        name?: string;
        orgId?: string[];
        /** @min 1 */
        pageNumber?: number;
        /**
         * @min 1
         * @max 100
         */
        pageSize?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<DomainProjectsPagedModel, HttpxErrBlock>({
        path: `/projects`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Creates a new project in an organisation
     *
     * @tags project
     * @name ProjectsCreate
     * @summary Create a project
     * @request POST:/projects
     * @secure
     */
    projectsCreate: (
      query: {
        /** Organisation ID */
        orgId: string;
      },
      body: DomainProjectCreateModel,
      params: RequestParams = {},
    ) =>
      this.request<DomainProjectModel, HttpxErrBlock>({
        path: `/projects`,
        method: "POST",
        query: query,
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Returns a single project by ID
     *
     * @tags project
     * @name ProjectsDetail
     * @summary Get a project
     * @request GET:/projects/{id}
     * @secure
     */
    projectsDetail: (id: string, params: RequestParams = {}) =>
      this.request<DomainProjectModel, HttpxErrBlock>({
        path: `/projects/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Soft deletes a project
     *
     * @tags project
     * @name ProjectsDelete
     * @summary Delete a project
     * @request DELETE:/projects/{id}
     * @secure
     */
    projectsDelete: (id: string, params: RequestParams = {}) =>
      this.request<void, HttpxErrBlock>({
        path: `/projects/${id}`,
        method: "DELETE",
        secure: true,
        ...params,
      }),

    /**
     * @description Updates a project's name and description
     *
     * @tags project
     * @name ProjectsPartialUpdate
     * @summary Update a project
     * @request PATCH:/projects/{id}
     * @secure
     */
    projectsPartialUpdate: (
      id: string,
      body: DomainProjectUpdateModel,
      params: RequestParams = {},
    ) =>
      this.request<DomainProjectModel, HttpxErrBlock>({
        path: `/projects/${id}`,
        method: "PATCH",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Changes a project's visibility (public/private)
     *
     * @tags project
     * @name VisibilityPartialUpdate
     * @summary Update project visibility
     * @request PATCH:/projects/{id}/visibility
     * @secure
     */
    visibilityPartialUpdate: (
      id: string,
      body: DomainProjectVisibilityModel,
      params: RequestParams = {},
    ) =>
      this.request<DomainProjectModel, HttpxErrBlock>({
        path: `/projects/${id}/visibility`,
        method: "PATCH",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  sprints = {
    /**
     * @description Returns paginated sprints in a project with optional filtering
     *
     * @tags sprint
     * @name SprintsList
     * @summary List sprints with pagination
     * @request GET:/sprints
     * @secure
     */
    sprintsList: (
      query?: {
        id?: string[];
        name?: string;
        /** @min 1 */
        pageNumber?: number;
        /**
         * @min 1
         * @max 100
         */
        pageSize?: number;
        projectId?: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<DomainSprintsPagedModel, HttpxErrBlock>({
        path: `/sprints`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Creates a new sprint in a project
     *
     * @tags sprint
     * @name SprintsCreate
     * @summary Create a sprint
     * @request POST:/sprints
     * @secure
     */
    sprintsCreate: (
      body: DomainSprintCreateModel,
      params: RequestParams = {},
    ) =>
      this.request<DomainSprintModel, HttpxErrBlock>({
        path: `/sprints`,
        method: "POST",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Returns a single sprint by ID
     *
     * @tags sprint
     * @name SprintsDetail
     * @summary Get a sprint
     * @request GET:/sprints/{sprintId}
     * @secure
     */
    sprintsDetail: (sprintId: string, params: RequestParams = {}) =>
      this.request<DomainSprintModel, HttpxErrBlock>({
        path: `/sprints/${sprintId}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Updates sprint details
     *
     * @tags sprint
     * @name SprintsPartialUpdate
     * @summary Update a sprint
     * @request PATCH:/sprints/{sprintId}
     * @secure
     */
    sprintsPartialUpdate: (
      sprintId: string,
      body: DomainSprintUpdateModel,
      params: RequestParams = {},
    ) =>
      this.request<DomainSprintModel, HttpxErrBlock>({
        path: `/sprints/${sprintId}`,
        method: "PATCH",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Transitions a sprint to completed status
     *
     * @tags sprint
     * @name CompletedCreate
     * @summary Complete a sprint
     * @request POST:/sprints/{sprintId}/completed
     * @secure
     */
    completedCreate: (sprintId: string, params: RequestParams = {}) =>
      this.request<DomainSprintModel, HttpxErrBlock>({
        path: `/sprints/${sprintId}/completed`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Transitions a sprint to active status
     *
     * @tags sprint
     * @name StartCreate
     * @summary Start a sprint
     * @request POST:/sprints/{sprintId}/start
     * @secure
     */
    startCreate: (sprintId: string, params: RequestParams = {}) =>
      this.request<DomainSprintModel, HttpxErrBlock>({
        path: `/sprints/${sprintId}/start`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),
  };
  tickets = {
    /**
     * @description Returns paginated tickets for a project, optionally filtered by sprint or board
     *
     * @tags ticket
     * @name TicketsList
     * @summary List tickets with pagination
     * @request GET:/tickets
     * @secure
     */
    ticketsList: (
      query?: {
        boardId?: string[];
        id?: string[];
        /** @min 1 */
        pageNumber?: number;
        /**
         * @min 1
         * @max 100
         */
        pageSize?: number;
        projectId?: string[];
        sprintId?: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<DomainTicketsPagedModel, HttpxErrBlock>({
        path: `/tickets`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Creates a new ticket in a project
     *
     * @tags ticket
     * @name TicketsCreate
     * @summary Create a ticket
     * @request POST:/tickets
     * @secure
     */
    ticketsCreate: (
      query: {
        /** Project ID */
        projectId: string;
      },
      body: DomainTicketCreateModel,
      params: RequestParams = {},
    ) =>
      this.request<DomainTicketModel, HttpxErrBlock>({
        path: `/tickets`,
        method: "POST",
        query: query,
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Returns a single ticket by ID
     *
     * @tags ticket
     * @name TicketsDetail
     * @summary Get a ticket
     * @request GET:/tickets/{ticketId}
     * @secure
     */
    ticketsDetail: (ticketId: string, params: RequestParams = {}) =>
      this.request<DomainTicketModel, HttpxErrBlock>({
        path: `/tickets/${ticketId}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Soft-deletes a ticket by ID
     *
     * @tags ticket
     * @name TicketsDelete
     * @summary Delete a ticket
     * @request DELETE:/tickets/{ticketId}
     * @secure
     */
    ticketsDelete: (ticketId: string, params: RequestParams = {}) =>
      this.request<void, HttpxErrBlock>({
        path: `/tickets/${ticketId}`,
        method: "DELETE",
        secure: true,
        ...params,
      }),

    /**
     * @description Updates ticket details (title, description, priority, type, etc.)
     *
     * @tags ticket
     * @name TicketsPartialUpdate
     * @summary Update a ticket
     * @request PATCH:/tickets/{ticketId}
     * @secure
     */
    ticketsPartialUpdate: (
      ticketId: string,
      body: DomainTicketUpdateModel,
      params: RequestParams = {},
    ) =>
      this.request<DomainTicketModel, HttpxErrBlock>({
        path: `/tickets/${ticketId}`,
        method: "PATCH",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Moves a ticket to a specific board column
     *
     * @tags ticket
     * @name MoveBoardColumnPartialUpdate
     * @summary Move ticket to board column
     * @request PATCH:/tickets/{ticketId}/move-board-column
     * @secure
     */
    moveBoardColumnPartialUpdate: (
      ticketId: string,
      body: DomainTicketBoardMoveModel,
      params: RequestParams = {},
    ) =>
      this.request<DomainTicketModel, HttpxErrBlock>({
        path: `/tickets/${ticketId}/move-board-column`,
        method: "PATCH",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Moves a ticket to a specific board and column
     *
     * @tags ticket
     * @name MoveToBoardPartialUpdate
     * @summary Move ticket to board column
     * @request PATCH:/tickets/{ticketId}/move-to-board
     * @secure
     */
    moveToBoardPartialUpdate: (
      ticketId: string,
      body: DomainTicketBoardMoveModel,
      params: RequestParams = {},
    ) =>
      this.request<DomainTicketModel, HttpxErrBlock>({
        path: `/tickets/${ticketId}/move-to-board`,
        method: "PATCH",
        body: body,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Moves a ticket to a specific sprint
     *
     * @tags ticket
     * @name MoveToSprintPartialUpdate
     * @summary Move ticket to sprint
     * @request PATCH:/tickets/{ticketId}/move-to-sprint
     * @secure
     */
    moveToSprintPartialUpdate: (
      ticketId: string,
      query: {
        /** Sprint ID */
        sprintId: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<DomainTicketModel, HttpxErrBlock>({
        path: `/tickets/${ticketId}/move-to-sprint`,
        method: "PATCH",
        query: query,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  users = {
    /**
     * @description Returns the authenticated user's profile
     *
     * @tags user
     * @name GetUsers
     * @summary Get current user
     * @request GET:/users/me
     * @secure
     */
    getUsers: (params: RequestParams = {}) =>
      this.request<DomainUserModel, HttpxErrBlock>({
        path: `/users/me`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),
  };
}
