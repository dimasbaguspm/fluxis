import { useListOrgs } from "@/hooks/use-api";
import { Icon, Loader, Text } from "@versaur/react/primitive";
import { ComboboxInput } from "@versaur/react/forms";
import { useCallback, useMemo, useState } from "react";
import { debounce } from "radash";
import { SearchIcon } from "@versaur/icons";
import { Controller } from "react-hook-form";
import type { Control, FieldValues, Path } from "react-hook-form";

interface SelectOrganisationsInputProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
  label?: string;
  placeholder?: string;
  helper?: string;
  required?: boolean;
  disabled?: boolean;
  multiple?: boolean;
}

export const SelectOrganisationsInput = <T extends FieldValues>({
  control,
  name,
  label = "Organisation",
  placeholder = "Search organisations...",
  helper,
  required,
  disabled,
  multiple = false,
}: SelectOrganisationsInputProps<T>) => {
  const [inputValue, setInputValue] = useState("");
  const [searchTerm, setSearchTerm] = useState("");

  // Debounced search handler
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
      setInputValue(value); // Immediate UI update
      debouncedSearch(value); // Debounced API call
    },
    [debouncedSearch],
  );

  const query = useMemo(() => ({ pageSize: 15, name: [searchTerm] }), [searchTerm]);

  const [orgs, err, { isPending }] = useListOrgs(query, {
    enabled: searchTerm.length > 0,
  });

  const [selectedOrg, setSelectedOrg] = useState<any>(null);

  const handleSelectionChange = useCallback(
    (newValue: string | string[] | null, onChange: (value: any) => void) => {
      if (multiple) {
        const values = Array.isArray(newValue) ? newValue : newValue ? [newValue] : [];
        onChange(values);
      } else {
        const value = Array.isArray(newValue) ? newValue[0] : newValue;
        onChange(value || null);
        // Track selected org for displaying label
        if (value && orgs?.items) {
          const selected = orgs.items.find((org) => org.id === value);
          if (selected) {
            setSelectedOrg(selected);
          }
        }
      }
    },
    [multiple, orgs],
  );

  const getDisplayOptions = useCallback(() => {
    const searchResults = orgs?.items || [];
    // Include selected org even if not in current search results
    if (selectedOrg && !searchResults.find((org) => org.id === selectedOrg.id)) {
      return [selectedOrg, ...searchResults];
    }
    return searchResults;
  }, [orgs, selectedOrg]);

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
                ? "Select organisations"
                : selectedOrg
                  ? selectedOrg.name
                  : "Select organisation"}
            </ComboboxInput.Button>
            <ComboboxInput.Container
              variant="list"
              search={
                <ComboboxInput.Search
                  name="org-search"
                  value={inputValue}
                  onChange={handleSearchChange}
                  placeholder="Search organisations..."
                />
              }
            >
              {inputValue.length === 0 ? (
                <Text as="span">Start typing to search organisations</Text>
              ) : isPending ? (
                <Loader />
              ) : err ? (
                <Text as="span">Error loading organisations</Text>
              ) : getDisplayOptions().length === 0 ? (
                <Text as="span">No organisations found</Text>
              ) : (
                getDisplayOptions().map((org) => (
                  <ComboboxInput.Option key={org.id} value={org.id}>
                    {org.name}
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
