import { useGetTicket, useListTickets } from "@/hooks/use-api";
import { vM2 } from "@versaur/core/utilities";
import { SearchIcon } from "@versaur/icons";
import { ComboboxInput } from "@versaur/react/forms";
import { Icon, Loader, Text } from "@versaur/react/primitive";
import { debounce } from "radash";
import { useCallback, useMemo, useState } from "react";
import type { Control, FieldValues, Path } from "react-hook-form";
import { Controller, useWatch } from "react-hook-form";

interface SelectTicketsInputProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
  label?: string;
  placeholder?: string;
  helper?: string;
  required?: boolean;
  disabled?: boolean;
  multiple?: boolean;
  projectId?: string;
  sprintId?: string;
  boardId?: string;
}

export const SelectTicketsInput = <T extends FieldValues>({
  control,
  name,
  label = "Ticket",
  placeholder = "Search tickets...",
  helper,
  required,
  disabled,
  multiple = false,
  projectId,
  sprintId,
  boardId,
}: SelectTicketsInputProps<T>) => {
  const [inputValue, setInputValue] = useState("");
  const [, setSearchTerm] = useState("");

  const val = useWatch({ control, name });

  const debouncedSearch = useMemo(
    () =>
      debounce({ delay: 300 }, (term: string) => {
        setSearchTerm(term);
      }),
    [],
  );

  const handleSearchChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const value = e.currentTarget.value;
      setInputValue(value);
      debouncedSearch(value);
    },
    [debouncedSearch],
  );

  const query = useMemo(
    () => ({
      pageSize: 15,
      ...(projectId ? { projectId: [projectId] } : {}),
      ...(sprintId ? { sprintId: [sprintId] } : {}),
      ...(boardId ? { boardId: [boardId] } : {}),
    }),
    [projectId, sprintId, boardId],
  );

  const [tickets, err, { isPending }] = useListTickets(query);
  const [ticket] = useGetTicket(val);

  const selectedTicket = useMemo(() => {
    if (!val) return null;

    if (tickets?.items) {
      const found = tickets.items.find((t) => t.id === val);
      if (found) return found;
    }
    return ticket;
  }, [val, tickets?.items, ticket]);

  const handleSelectionChange = useCallback(
    (newValue: string | string[] | null, onChange: (value: any) => void) => {
      if (multiple) {
        const values = Array.isArray(newValue) ? newValue : newValue ? [newValue] : [];
        onChange(values);
      } else {
        const value = Array.isArray(newValue) ? newValue[0] : newValue;
        onChange(value || null);
      }
    },
    [multiple],
  );

  const getDisplayOptions = useCallback(() => {
    return tickets?.items || [];
  }, [tickets]);

  const getTicketLabel = useCallback(
    (ticket: any) => `${ticket.key || ticket.id} - ${ticket.title}`,
    [],
  );

  return (
    <Controller
      name={name}
      control={control}
      rules={{ required: required ? `${label} is required` : undefined }}
      render={({ field, fieldState: { error } }) => {
        const listboxContent = (
          <>
            <ComboboxInput.Button>
              {multiple
                ? "Select tickets"
                : selectedTicket
                  ? getTicketLabel(selectedTicket)
                  : "Select ticket"}
            </ComboboxInput.Button>
            <ComboboxInput.Container
              variant="list"
              search={
                <ComboboxInput.Search
                  name="ticket-search"
                  value={inputValue}
                  onChange={handleSearchChange}
                  placeholder="Search tickets..."
                />
              }
            >
              {isPending ? (
                <Loader type="bar" />
              ) : err ? (
                <Text as="span" className={vM2}>
                  Error loading tickets
                </Text>
              ) : getDisplayOptions().length === 0 ? (
                <Text as="span" className={vM2}>
                  No tickets found
                </Text>
              ) : (
                getDisplayOptions().map((ticket) => (
                  <ComboboxInput.Option key={ticket.id} value={ticket.id}>
                    {getTicketLabel(ticket)}
                  </ComboboxInput.Option>
                ))
              )}
            </ComboboxInput.Container>
            {multiple && <ComboboxInput.SelectionChips />}
          </>
        );

        return multiple ? (
          <ComboboxInput
            multiple
            value={(field.value as string[]) || []}
            onChange={(newValue) => handleSelectionChange(newValue, field.onChange)}
            label={label}
            placeholder={placeholder}
            helper={helper}
            required={required}
            disabled={disabled}
            invalid={!!error}
            error={error?.message}
            iconLeft={<Icon as={SearchIcon} />}
          >
            {listboxContent}
          </ComboboxInput>
        ) : (
          <ComboboxInput
            multiple={false}
            value={(field.value as string) || ""}
            onChange={(newValue) => handleSelectionChange(newValue, field.onChange)}
            label={label}
            placeholder={placeholder}
            helper={helper}
            required={required}
            disabled={disabled}
            invalid={!!error}
            error={error?.message}
            iconLeft={<Icon as={SearchIcon} />}
          >
            {listboxContent}
          </ComboboxInput>
        );
      }}
    />
  );
};
