import { useListBoards } from "@/hooks/use-api";
import { Icon, Loader, Text } from "@versaur/react/primitive";
import { ComboboxInput } from "@versaur/react/forms";
import { useCallback, useMemo, useState } from "react";
import { debounce } from "radash";
import { SearchIcon } from "@versaur/icons";
import { Controller } from "react-hook-form";
import type { Control, FieldValues, Path } from "react-hook-form";

interface SelectBoardsInputProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
  label?: string;
  placeholder?: string;
  helper?: string;
  required?: boolean;
  disabled?: boolean;
  multiple?: boolean;
  sprintId?: string;
}

export const SelectBoardsInput = <T extends FieldValues>({
  control,
  name,
  label = "Board",
  placeholder = "Search boards...",
  helper,
  required,
  disabled,
  multiple = false,
  sprintId,
}: SelectBoardsInputProps<T>) => {
  const [inputValue, setInputValue] = useState("");
  const [searchTerm, setSearchTerm] = useState("");

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
    () => ({ pageSize: 15, name: searchTerm, ...(sprintId ? { sprintId: [sprintId] } : {}) }),
    [searchTerm, sprintId],
  );

  const [boards, err, { isPending }] = useListBoards(query, {
    enabled: searchTerm.length > 0,
  });

  const [selectedBoard, setSelectedBoard] = useState<any>(null);

  const handleSelectionChange = useCallback(
    (newValue: string | string[] | null, onChange: (value: any) => void) => {
      if (multiple) {
        const values = Array.isArray(newValue) ? newValue : newValue ? [newValue] : [];
        onChange(values);
      } else {
        const value = Array.isArray(newValue) ? newValue[0] : newValue;
        onChange(value || null);
        if (value && boards?.items) {
          const selected = boards.items.find((board) => board.id === value);
          if (selected) {
            setSelectedBoard(selected);
          }
        }
      }
    },
    [multiple, boards],
  );

  const getDisplayOptions = useCallback(() => {
    const searchResults = boards?.items || [];
    if (selectedBoard && !searchResults.find((board) => board.id === selectedBoard.id)) {
      return [selectedBoard, ...searchResults];
    }
    return searchResults;
  }, [boards, selectedBoard]);

  return (
    <Controller
      name={name}
      control={control}
      rules={{ required: required ? `${label} is required` : undefined }}
      render={({ field, fieldState: { error } }) => {
        const listboxContent = (
          <>
            <ComboboxInput.Button>
              {multiple ? "Select boards" : selectedBoard ? selectedBoard.name : "Select board"}
            </ComboboxInput.Button>
            <ComboboxInput.Container
              variant="list"
              search={
                <ComboboxInput.Search
                  name="board-search"
                  value={inputValue}
                  onChange={handleSearchChange}
                  placeholder="Search boards..."
                />
              }
            >
              {inputValue.length === 0 ? (
                <Text as="span">Start typing to search boards</Text>
              ) : isPending ? (
                <Loader />
              ) : err ? (
                <Text as="span">Error loading boards</Text>
              ) : getDisplayOptions().length === 0 ? (
                <Text as="span">No boards found</Text>
              ) : (
                getDisplayOptions().map((board) => (
                  <ComboboxInput.Option key={board.id} value={board.id}>
                    {board.name}
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
