import type { DomainAuthModel, DomainUserModel } from "@interfaces/openapi.generated";

export type User = DomainUserModel;
export type Auth = DomainAuthModel;

export interface Session {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
}
