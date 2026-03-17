import { useListSprints } from "@/hooks/use-api";
import { Icon, Loader, Text } from "@versaur/react/primitive";
import { ComboboxInput } from "@versaur/react/forms";
import { useCallback, useMemo, useState } from "react";
import { debounce } from "radash";
import { SearchIcon } from "@versaur/icons";
import { Controller } from "react-hook-form";
import type { Control, FieldValues, Path } from "react-hook-form";

interface SelectSprintsInputProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
  label?: string;
  placeholder?: string;
  helper?: string;
  required?: boolean;
  disabled?: boolean;
  multiple?: boolean;
  projectId?: string;
}

export const SelectSprintsInput = <T extends FieldValues>({
  control,
  name,
  label = "Sprint",
  placeholder = "Search sprints...",
  helper,
  required,
  disabled,
  multiple = false,
  projectId,
}: SelectSprintsInputProps<T>) => {
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
    () => ({ pageSize: 15, name: searchTerm, ...(projectId ? { projectId: [projectId] } : {}) }),
    [searchTerm, projectId],
  );

  const [sprints, err, { isPending }] = useListSprints(query, {
    enabled: searchTerm.length > 0,
  });

  const [selectedSprint, setSelectedSprint] = useState<any>(null);

  const handleSelectionChange = useCallback(
    (newValue: string | string[] | null, onChange: (value: any) => void) => {
      if (multiple) {
        const values = Array.isArray(newValue) ? newValue : newValue ? [newValue] : [];
        onChange(values);
      } else {
        const value = Array.isArray(newValue) ? newValue[0] : newValue;
        onChange(value || null);
        if (value && sprints?.items) {
          const selected = sprints.items.find((sprint) => sprint.id === value);
          if (selected) {
            setSelectedSprint(selected);
          }
        }
      }
    },
    [multiple, sprints],
  );

  const getDisplayOptions = useCallback(() => {
    const searchResults = sprints?.items || [];
    if (selectedSprint && !searchResults.find((sprint) => sprint.id === selectedSprint.id)) {
      return [selectedSprint, ...searchResults];
    }
    return searchResults;
  }, [sprints, selectedSprint]);

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
                ? "Select sprints"
                : selectedSprint
                  ? selectedSprint.name
                  : "Select sprint"}
            </ComboboxInput.Button>
            <ComboboxInput.Container
              variant="list"
              search={
                <ComboboxInput.Search
                  name="sprint-search"
                  value={inputValue}
                  onChange={handleSearchChange}
                  placeholder="Search sprints..."
                />
              }
            >
              {inputValue.length === 0 ? (
                <Text as="span">Start typing to search sprints</Text>
              ) : isPending ? (
                <Loader />
              ) : err ? (
                <Text as="span">Error loading sprints</Text>
              ) : getDisplayOptions().length === 0 ? (
                <Text as="span">No sprints found</Text>
              ) : (
                getDisplayOptions().map((sprint) => (
                  <ComboboxInput.Option key={sprint.id} value={sprint.id}>
                    {sprint.name}
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
