-- mermaid-filter.lua
local system = require 'pandoc.system'

-- Helper to check if file exists
local function file_exists(name)
  local f = io.open(name, "r")
  if f ~= nil then io.close(f) return true else return false end
end

-- Helper to compute sha1 (requires openssl installed usually, or we can use pandoc built-in)
-- Pandoc 2.0+ has pandoc.sha1
function CodeBlock(cb)
  if cb.classes:includes('mermaid') then
    -- Generate hash for filename
    local hash = pandoc.sha1(cb.text)
    local img_filename = "mermaid-" .. hash .. ".png"
    
    -- Compile if not exists
    if not file_exists(img_filename) then
      -- Write mermaid code to temp file
      local cmd_input = "mermaid-" .. hash .. ".mmd"
      local f = io.open(cmd_input, "w")
      f:write(cb.text)
      f:close()
      
      -- Execute mmdc
      -- Redirect stderr to stdout so we can see errors in the log if needed
      local cmd = string.format("mmdc -i %s -o %s -b transparent 2>&1", cmd_input, img_filename)
      -- print("Executing: " .. cmd) -- Debug print removed
      local handle = io.popen(cmd)
      local result = handle:read("*a")
      local success, exit_reason, exit_code = handle:close()
      
      -- Handle Lua 5.1/Pandoc vs 5.3+ return values
      if type(success) == "nil" then
           -- For some Lua versions, io.popen:close returns exit code directly
      end
      
      if not file_exists(img_filename) then
        print("Failed to generate mermaid diagram for block " .. hash)
        print("Output: " .. (result or ""))
        return cb
      end
      
      -- cleanup temp input file only on success
      os.remove(cmd_input)
    end
    
    -- Return Image
    return pandoc.Image({pandoc.Str("Mermaid Diagram")}, img_filename, "generated diagram")
  end
end
