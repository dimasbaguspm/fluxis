import { useListOrgs } from "@/hooks/use-api";
import { Icon } from "@versaur/react/primitive";
import { ComboboxInput } from "@versaur/react/forms";
import { useCallback, useState } from "react";
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
}

export const SelectOrganisationsInput = <T extends FieldValues>({
  control,
  name,
  label = "Organisation",
  placeholder = "Search organisations...",
  helper,
  required,
  disabled,
}: SelectOrganisationsInputProps<T>) => {
  const [searchTerm, setSearchTerm] = useState("");
  const [isOpen, setIsOpen] = useState(false);

  // Fetch organisations based on search term
  const [orgs] = useListOrgs(
    isOpen || searchTerm
      ? { name: searchTerm ? [searchTerm] : undefined, pageSize: 20 }
      : undefined,
  );

  // Debounced search handler
  const debouncedSearch = useCallback(
    debounce({ delay: 300 }, (term: string) => {
      setSearchTerm(term);
    }),
    [],
  );

  const handleSelectionChange = useCallback(
    (newValues: string[], onChange: (value: string) => void) => {
      if (newValues.length > 0) {
        onChange(newValues[0]);
        setIsOpen(false);
      }
    },
    [],
  );

  return (
    <Controller
      name={name}
      control={control}
      rules={{ required: required ? `${label} is required` : undefined }}
      render={({ field, fieldState: { error } }) => (
        <ComboboxInput
          value={field.value ? [field.value] : []}
          onChange={(newValues) => handleSelectionChange(newValues, field.onChange)}
          label={label}
          placeholder={placeholder}
          helper={helper}
          required={required}
          disabled={disabled}
          invalid={!!error}
          error={error?.message}
          iconLeft={<Icon as={SearchIcon} />}
        >
          <ComboboxInput.Button />
          <ComboboxInput.Listbox>
            {orgs?.items && orgs.items.length > 0 ? (
              orgs.items.map((org) => (
                <ComboboxInput.Option key={org.id} value={org.id}>
                  {org.name}
                </ComboboxInput.Option>
              ))
            ) : (
              <div style={{ padding: "8px", textAlign: "center", color: "#999", fontSize: "12px" }}>
                {searchTerm ? "No organisations found" : isOpen ? "Type to search" : ""}
              </div>
            )}
          </ComboboxInput.Listbox>
        </ComboboxInput>
      )}
    />
  );
};
