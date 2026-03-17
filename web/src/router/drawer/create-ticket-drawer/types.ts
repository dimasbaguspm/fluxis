export interface CreateTicketFormInputs {
  projectId: string;
  title: string;
  type: string;
  priority: string;
  description?: string;
  sprintId?: string;
  storyPoints?: number;
  dueDate?: string;
}
