using HTTPClient
using Gumbo

function getPages(url, pages)
  page_array = [bytestring(get(string(url,1)).body)]
  for p in 2:pages
	push!(page_array, bytestring(get(string(url,p)).body))
  end
  return page_array
end

function parsePages(pages)
  bodies = [parsehtml(pages[1]).root[2]]
  for p in pages[2:end]
	push!(bodies, parsehtml(p).root[2])
  end
end


