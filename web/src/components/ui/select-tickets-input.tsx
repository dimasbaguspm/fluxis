import { useListTickets } from "@/hooks/use-api";
import { Icon, Loader, Text } from "@versaur/react/primitive";
import { ComboboxInput } from "@versaur/react/forms";
import { useCallback, useMemo, useState } from "react";
import { debounce } from "radash";
import { SearchIcon } from "@versaur/icons";
import { Controller } from "react-hook-form";
import type { Control, FieldValues, Path } from "react-hook-form";

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
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedTicket, setSelectedTicket] = useState<any>(null);

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

  const [tickets, err, { isPending }] = useListTickets(query, {
    enabled: searchTerm.length > 0,
  });

  const handleSelectionChange = useCallback(
    (newValue: string | string[] | null, onChange: (value: any) => void) => {
      if (multiple) {
        const values = Array.isArray(newValue) ? newValue : newValue ? [newValue] : [];
        onChange(values);
      } else {
        const value = Array.isArray(newValue) ? newValue[0] : newValue;
        onChange(value || null);
        if (value && tickets?.items) {
          const selected = tickets.items.find((ticket) => ticket.id === value);
          if (selected) {
            setSelectedTicket(selected);
          }
        }
      }
    },
    [multiple, tickets],
  );

  const getDisplayOptions = useCallback(() => {
    const searchResults = tickets?.items || [];
    if (selectedTicket && !searchResults.find((ticket) => ticket.id === selectedTicket.id)) {
      return [selectedTicket, ...searchResults];
    }
    return searchResults;
  }, [tickets?.items, selectedTicket]);

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
              {inputValue.length === 0 ? (
                <Text as="span">Start typing to search tickets</Text>
              ) : isPending ? (
                <Loader />
              ) : err ? (
                <Text as="span">Error loading tickets</Text>
              ) : getDisplayOptions().length === 0 ? (
                <Text as="span">No tickets found</Text>
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
