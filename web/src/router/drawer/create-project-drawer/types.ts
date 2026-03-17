export interface CreateProjectFormInputs {
  orgId: string;
  name: string;
  key: string;
  visibility: "public" | "private";
  description?: string;
}
