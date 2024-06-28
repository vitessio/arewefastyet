/*
Copyright 2023 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import React, { useEffect, useRef, useState, ReactNode } from "react";
import useClickOutside from "../hooks/useClickOutside";
import Icon from "../common/Icon";
import { twMerge } from "tailwind-merge";
import { Button } from "@/components/ui/button";

interface ContainerProps {
  className?: string;
  defaultIndex?: number;
  placeholder?: string;
  onChange: (value: { value: string }) => void;
  children: ReactNode;
}

interface OptionProps {
  className?: string;
  isSelected?: boolean;
  onClick?: () => void;
  children: ReactNode;
}

/**
 * Dropdown container component that manages the state and rendering of the dropdown menu.
 *
 * @param {Object} props - The properties for the Container component.
 * @param {string} [props.className] - Additional class names for styling.
 * @param {number} [props.defaultIndex=0] - The default selected index.
 * @param {string} [props.placeholder] - The placeholder text when no item is selected.
 * @param {(value: { value: string }) => void} props.onChange - Callback function called when the selected item changes.
 * @param {React.ReactNode} props.children - The dropdown options as children.
 * @returns {JSX.Element} The rendered Dropdown component.
 */

function Container({
  className,
  defaultIndex = 0,
  placeholder,
  onChange,
  children,
}: ContainerProps): JSX.Element {
  const [open, setOpen] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(defaultIndex);

  const ref = useRef<HTMLDivElement>(null);

  let items =
    React.Children.map(
      children,
      (child) => `${(child as React.ReactElement).props.children}`
    ) || [];

  useClickOutside(ref as React.MutableRefObject<HTMLElement>, () =>
    setOpen(false)
  );

  function handleOptionClick(index: number) {
    setOpen(false);
    setSelectedIndex(index);
  }

  useEffect(() => {
    if (selectedIndex < 0) setSelectedIndex(0);
  }, [selectedIndex]);

  useEffect(() => {
    onChange({ value: items[selectedIndex] });
  }, [selectedIndex]);

  return (
    <div ref={ref} className="relative flex gap-2">
      <Button
        variant="outline"
        onClick={() => setOpen(!open)}
        className={twMerge(
          "flex items-center h-full justify-center duration-inherit hover:bg-transparent",
          className
        )}
      >
        {`${items[selectedIndex] || placeholder}`}
        <Icon
          icon="expand_more"
          className={twMerge(
            "pt-[3px] ml-2 scale-150 duration-300",
            open && "rotate-180"
          )}
        />
      </Button>

      <>
        <div
          className={twMerge(
            "flex flex-col left-1/2 -translate-x-1/2 z-10 duration-300 absolute top-full"
          )}
          style={{
            clipPath: !open
              ? "polygon(0% 0%, 0% 0%, 100% 0%, 100% 0%)"
              : "polygon(0% 0%, 0% 100%, 100% 100%, 100% 0%)",
          }}
        >
          {React.Children.map(children, (child, index) => {
            if (React.isValidElement(child) && child.type === Option)
              return React.cloneElement(child, {
                index,
                onClick: () => handleOptionClick(index),
                isSelected: index === selectedIndex,
              } as Partial<OptionProps>);
          })}
        </div>
      </>
    </div>
  );
}

/**
 * Dropdown option component that renders an individual dropdown item.
 *
 * @param {Object} props - The properties for the Option component.
 * @param {string} [props.className] - Additional class names for styling.
 * @param {boolean} [props.isSelected] - Whether the option is currently selected.
 * @param {() => void} [props.onClick] - Click handler for the option.
 * @param {React.ReactNode} props.children - The content of the option.
 * @returns {JSX.Element} The rendered Option component.
 */

function Option({
  className,
  isSelected,
  onClick,
  children,
}: OptionProps): JSX.Element {
  return (
    <Button
      className={twMerge(
        "dark:text-white text-black",
        className,
        isSelected && "bg-primary"
      )}
      onClick={onClick}
    >
      {children}
    </Button>
  );
}

const Dropdown = { Container, Option };

export default Dropdown;
