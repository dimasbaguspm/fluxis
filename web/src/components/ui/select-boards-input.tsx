import { useGetBoard, useListBoards } from "@/hooks/use-api";
import { vM2 } from "@versaur/core/utilities";
import { SearchIcon } from "@versaur/icons";
import { ComboboxInput } from "@versaur/react/forms";
import { Icon, Loader, Text } from "@versaur/react/primitive";
import { debounce } from "radash";
import { useCallback, useMemo, useState } from "react";
import type { Control, FieldValues, Path } from "react-hook-form";
import { Controller, useWatch } from "react-hook-form";

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
    () => ({ pageSize: 15, name: searchTerm, ...(sprintId ? { sprintId: [sprintId] } : {}) }),
    [searchTerm, sprintId],
  );

  const [boards, err, { isPending }] = useListBoards(query);
  const [board] = useGetBoard(val);

  const selectedBoard = useMemo(() => {
    if (!val) return null;

    if (boards?.items) {
      const found = boards.items.find((b) => b.id === val);
      if (found) return found;
    }

    return board;
  }, [val, boards?.items, board]);

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
    return boards?.items || [];
  }, [boards]);

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
              {isPending ? (
                <Loader type="bar" />
              ) : err ? (
                <Text as="span" className={vM2}>
                  Error loading boards
                </Text>
              ) : getDisplayOptions().length === 0 ? (
                <Text as="span" className={vM2}>
                  No boards found
                </Text>
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
