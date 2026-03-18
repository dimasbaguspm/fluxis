import { useListSprints, useGetSprint } from "@/hooks/use-api";
import { Icon, Loader, Text } from "@versaur/react/primitive";
import { ComboboxInput } from "@versaur/react/forms";
import { useCallback, useMemo, useState } from "react";
import { debounce } from "radash";
import { SearchIcon } from "@versaur/icons";
import { Controller, useWatch } from "react-hook-form";
import type { Control, FieldValues, Path } from "react-hook-form";
import { vM2 } from "@versaur/core/utilities";

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
    () => ({ pageSize: 15, name: searchTerm, ...(projectId ? { projectId: [projectId] } : {}) }),
    [searchTerm, projectId],
  );

  const [sprints, err, { isPending }] = useListSprints(query);
  const [sprint] = useGetSprint(val);

  const selectedSprint = useMemo(() => {
    if (!val) return null;

    if (sprints?.items) {
      const found = sprints.items.find((s) => s.id === val);
      if (found) return found;
    }

    return sprint;
  }, [val, sprints?.items, sprint]);

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
    return sprints?.items || [];
  }, [sprints]);

  return (
    <Controller
      name={name}
      control={control}
      rules={{ required: required ? `${label} is required` : undefined }}
      render={({ field, fieldState: { error } }) => {
        const listboxContent = (
          <>
            <ComboboxInput.Button>
              {multiple ? "Select sprints" : selectedSprint ? selectedSprint.name : "Select sprint"}
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
              {isPending ? (
                <Loader type="bar" />
              ) : err ? (
                <Text as="span" className={vM2}>
                  Error loading sprints
                </Text>
              ) : getDisplayOptions().length === 0 ? (
                <Text as="span" className={vM2}>
                  No sprints found
                </Text>
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
