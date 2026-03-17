import type {
  DomainAuthModel,
  DomainBoardColumnModel,
  DomainBoardModel,
  DomainOrganisationMemberModel,
  DomainOrganisationModel,
  DomainProjectModel,
  DomainSprintModel,
  DomainTicketModel,
  DomainUserModel,
} from "@/interfaces/openapi.generated";

/**
 * Typed notification event with entity type, action, and payload
 */
export interface NotifierEvent {
  type: string;
  action: string;
  payload:
    | DomainAuthModel
    | DomainUserModel
    | DomainOrganisationModel
    | DomainOrganisationMemberModel
    | DomainProjectModel
    | DomainSprintModel
    | DomainBoardModel
    | DomainBoardColumnModel
    | DomainTicketModel
    | unknown;
}

export interface NotifierContextType {
  event: NotifierEvent | null;
  isConnected: boolean;
}
