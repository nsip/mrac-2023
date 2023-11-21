require 'json'
require 'csv'

# Validation of MRAC JSON-LD output.
# Presupposes that it is run in a folder containing two folders: the JSON-LD
# output for MRAC (asn-json-ld), and the ASN JSON output for MRAC that it derives
# from (asn-json), both as *.json files. Further presupposes that the
# ScOT JSON-LD file has been downloaded from https://scot.curriculum.edu.au, as
# scot.jsonld, in the asn-json-ld folder.
#
# Validation checks:
# * Whether there are any duplicate IDs in the JSON tree
# * whether all JSON keys are namespaced as JSON-LD keys;
# * whether the gem:hasChild/gem:isChildOf pairs match;
# * whether cross-references to other parts of MRAC are pointing to the right kind
# of referent;
# * whether ScOT cross-references are legal;
# * whether the year levels assigned to nodes are consistent with their
# identifier;
# * whether the year levels of a child node are a subset of the year levels of a parent node;
# * whether dc:educationLevel and esa:nominalYearLevel have consistent values;
# * whether the namespaced nodes in JSON have been preserved in JSON-LD.

JSON_FOLDER = 'asn-json'.freeze
JSON_LD_FOLDER = 'asn-json-ld'.freeze
SCOT_FILE = File.join(JSON_LD_FOLDER, 'scot.jsonld').freeze
SCOT_MAPPING_FILE = 'mapping.csv'.freeze

def scot_id_parse(json, ids)
  case json
  when Hash
    if (k = id(json)) &&
       %r{^http://vocabulary.curriculum.edu.au/scot/}.match?(k)
      ids[k] = true
    end
    json.each_value do |m|
      scot_id_parse(m, ids)
    end
  when Array
    json.each { |j| scot_id_parse(j, ids) }
  end
  ids
end

def asn_jsonld_id_parse(json, ids)
  if (k = id(json)) && json['asn:statementLabel'] # if no label, it's a root of tree, not a real ID
    ids[k] and @err.puts "#{@f}: id_parse: #{k} is duplicated beween #{@ids[k]} and #{@f}!"
    ids[k] = { file: @f.sub(%r{^.+/}, ''), label: json['asn:statementLabel'] }
  end
  json['gem:hasChild']&.each do |m|
    ids = asn_jsonld_id_parse(m, ids)
  end
  case json['asn:hasLevel']
  when Array
    json['asn:hasLevel']&.each do |m|
      ids = asn_jsonld_id_parse(m, ids)
    end
  when Hash
    ids = asn_jsonld_id_parse(json['asn:hasLevel'], ids)
  end
  ids
end

# parse not just IDs but also all predicates for each node
def asn_json_id_parse(json, ids)
  if (k = id(json))
    keys = json.keys.reject { |x| json[x].nil? }
    untitled = ['Cross-Curriculum Priority']
    # if !json['asn_statementLabel'] || untitled.include?(json['asn_statementLabel']) ||
    #   untitled.include?(json.dig('asn_statementLabel', 'literal'))
    # More trouble than it's worth
    if keys.include?('dcterms_description')
      keys = keys.reject { |x| x == 'dcterms_description' }
      keys.include?('dcterms_title') or keys << 'dcterms_title'
    end
    keys = keys.reject { |x| %w[esa_nominalYearLevel].include?(x) }
    # end
    ids[k] = { file: @f.sub(%r{^.+/}, ''), keys: keys.map { |x| x.sub('_', ':').sub('dcterms:', 'dc:') }
                                                     .select { |x| x.include?(':') }.sort }
    ids['children'] and ids[k] += 'gem:hasChild'
  end
  json['children']&.each do |m|
    ids = asn_json_id_parse(m, ids)
  end
  case json['asn_hasLevel']
  when Array
    json['asn_hasLevel']&.each do |m|
      ids = asn_json_id_parse(m, ids)
    end
  when Hash
    ids = asn_json_id_parse(json['asn_hasLevel'], ids)
  end
  ids
end

def child_parent(json)
  json['gem:hasChild']&.each do |j|
    j['gem:isChildOf']&.each do |k|
      id(json) or next
      (id(json) == id(k)) and next
      @err.puts "#{@f}: child_parent: #{id(j)} mismatches gem:hasChild and gem:isChildOf!"
    end
    child_parent(j)
  end
  json['asn:hasLevel']&.each do |j|
    #j["asn:statementNotation"] and
      #puts "#{json["asn:statementNotation"]} : #{json["asn:statementLabel"]} hasLevel #{j["asn:statementNotation"]} : #{j["asn:statementLabel"]}"
    j['gem:isChildOf']&.each do |k|
      id(json) or next
      (id(json) == id(k)) and next
      @err.puts "#{@f}: child_parent: #{id(j)} mismatches asn:hasLevel and gem:isChildOf!"
    end
    child_parent(j)
  end
end

def id(x)
  require 'debug'
  x.is_a?(Hash) or binding.b
  x['@id'] || x['@Id'] || x['id'] || x['Id']
end

def namespaced_keys(json, key)
  case json
  when Hash
    json.each_key do |k|
      namespaced_keys(json[k], k)
      key == '@context' and next
      !/\s/.match?(k) && /\S:\S|^@/.match?(k) and next
      @err.puts "#{@f}: namespaced_keys: #{k} in #{id(json)} is an illegal key!"
    end
  when Array
    json.each { |j| namespaced_keys(j, nil) }
  end
end

def exist_link(k, lbl)
  @ids[k] or
    @err.puts "#{@f}: link_types: #{k} destination of #{lbl} is not a defined key!"
  @ids[k]
end

def link_types(json)
  case json
  when Hash
    # require "debug"; binding.b if (@f == "la-English.json.json")
    link_skillembodied(json)
    link_cross_subject_reference(json)
    link_has_level(json)
    link_dc_relation(json)
    json.each_value do |m|
      link_types(m)
    end
  when Array
    json.each { |j| link_types(j) }
  end
rescue StandardError => e
  require 'debug'
  binding.b
end

def array_link_check(json, key)
  json[key] or return false
  unless json[key].is_a?(Array)
    @err.puts("#{@f}: link_types: #{key} under @id = #{json['@id']} is not an array!")
    return false
  end
  if json[key].empty?
    @err.puts("#{@f}: link_types: #{key} under @id = #{json['@id']} is empty!")
    return false
  end
  true
end

def link_skillembodied(json)
  array_link_check(json, 'asn:skillEmbodied') or return
  json['asn:skillEmbodied'].each do |k|
    exist_link(k['@id'], 'asn:skillEmbodied') or next
    unless /^gc/.match?(@ids[k['@id']][:file])
      @err.puts "#{@f}: link_types: #{k['@id']} destination of asn:skillEmbodied in #{json['@id']} is not a general capability!"
    end
  end
rescue StandardError => e
  require 'debug'
  binding.b
end

def link_cross_subject_reference(json)
  array_link_check(json, 'asn:crossSubjectReference') or return
  json['asn:crossSubjectReference'].each do |k|
    exist_link(k['@id'], 'asn:crossSubjectReference') or next
    /^ccp|^AA|^A_TSI|^S/.match?(@ids[k['@id']][:file]) or
      @err.puts "#{@f}: link_types: #{k['@id']} destination of asn:crossSubjectReference in #{json['@id']} is not a cross-curriculum priority!"
  end
rescue StandardError => e
  require 'debug'
  binding.b
end

def link_has_level(json)
  array_link_check(json, 'asn:hasLevel') or return
  json['asn:hasLevel'].each do |k|
    exist_link(k['@id'], 'asn:hasLevel') or next
    ['Level', 'Achievement Standard', 'Achievement Standard Component'].include?(@ids[k['@id']][:label]) or
      @err.puts "#{@f}: link_types: #{k['@id']} destination of asn:hasLevel in #{json['@id']} is #{@ids[k['@id']][:label]}, not an achievement level!"
  end
rescue StandardError => e
  require 'debug'
  binding.b
end

def link_dc_relation(json)
  array_link_check(json, 'dc:relation') or return
  json['dc:relation'].each do |k|
    exist_link(k['@id'], 'dc:relation') or next
    /^la/.match?(@ids[k['@id']][:file]) or
      @err.puts "#{@f}: link_types: #{k['@id']} destination of dc:relation in #{json['@id']} is not a curriculum statement!"
  end
rescue StandardError => e
  require 'debug'
  binding.b
end

def scot(json)
  case json
  when Hash
    (s = json['asn:conceptTerm']) and scot1(s)
    json.each_value do |m|
      scot(m)
    end
  when Array
    json.each { |j| scot(j) }
  end
end

def scot1(terms)
  ret = if terms.is_a?(Array) then terms.map { |x| x['@id'] }
        else
          [terms['@id']]
        end
  ret.each do |k|
    k.nil? and @err.puts "#{@f}: scot: #{terms} destination of asn:conceptTerm is undefined!"
    @scotids[k] or
      @err.puts "#{@f}: scot: #{k} destination of asn:conceptTerm is undefined!"
  end
end

def scot_mappings(json)
  case json
  when Hash
    if (s = json['asn:conceptTerm'])
      scot_mappings1(s, json['asn:statementNotation'])
    elsif @mappings[json['asn:statementNotation']]
      scot_mappings1([], json['asn:statementNotation'])
    end
    json.each do |_k, m|
      scot_mappings(m)
    end
  when Array
    json.each { |j| scot_mappings(j) }
  end
end

def scot_mappings1(terms, contentdesc)
  ret = if terms.is_a?(Array) then terms.map { |x| x['@id'] }
          elsif terms['@id']
          [terms['@id']]
          else nil
        end
  #@err.puts "#{@f}: scot mappings: #{contentdesc} SUCCESS" if @mappings[contentdesc] && @mappings[contentdesc] == ret
  @mappings[contentdesc] == ret and return
  error_report = terms.is_a?(String) ? terms : ret
  error_report.empty? and error_report = nil
  @err.puts "#{@f}: scot mappings: #{contentdesc} is meant to have ScOT terms #{@mappings[contentdesc] || "EMPTY"}, it has #{error_report || "EMPTY"}"
end

def year_level(json)
  consistent_nominal_yr_level(json)
  inferred_yr_level(json)
  inherit_yr_level(json)
end

def consistent_nominal_yr_level(json)
  if json['dc:educationLevel'] && json['esa:nominalYearLevel']
    lvl = extract_yr_levels(json['dc:educationLevel'])
    nomnl = extract_yr_levels(json['esa:nominalYearLevel'])
    unless subset?(lvl, nomnl)
      @err.puts "#{@f}: year levels: dc:educationLevel (#{lvl.join(', ')}) and esa:nominalYearLevel (#{nomnl.join(', ')})  in #{json['@id']} are inconsistent"
    end
  end
  json['gem:hasChild']&.each do |m|
    consistent_nominal_yr_level(m)
  end
  case json['asn:hasLevel']
  when Array
    json['asn:hasLevel']&.each do |m|
      consistent_nominal_yr_level(m)
    end
  when Hash
    consistent_nominal_yr_level(json['asn:hasLevel'])
  end
end

def inferred_yr_level(json)
  if a = json['asn:statementNotation'] and id = a['literal']
    /FY$/.match?(id) and id_yrs = ['F']
    m = /Y(\d+)$/.match?(id) and id_yrs = m[1].split
    lvl_yrs = extract_yr_levels(json['dc:educationLevel'])
    unless id_yrs.empty? || id_yrs.empty?
      subset?(id_yrs, lvl_yrs) or
        @err.puts "#{@f}: year levels: #{json['@id']} node year levels are inconsistent with identifier #{id} in #{json['@id']}"
    end
  end
  json['gem:hasChild']&.each do |m|
    inferred_yr_level(m)
  end
  case json['asn:hasLevel']
  when Array
    json['asn:hasLevel']&.each do |m|
      inferred_yr_level(m)
    end
  when Hash
    inferred_yr_level(json['asn:hasLevel'])
  end
end

def inherit_yr_level(json, parent_yrs = [])
  y = extract_yr_levels(json['dc:educationLevel'])
  if !(parent_yrs.empty? || y.empty?) && !subset?(parent_yrs, y)
    @err.puts "#{@f}: year levels: #{id(json)} node year levels (#{y.join(', ')}) are not subset of parent node year levels (#{parent_yrs.join(', ')})"
  end
  json['gem:hasChild']&.each do |m|
    inherit_yr_level(m, y)
  end
  case json['asn:hasLevel']
  when Array
    json['asn:hasLevel']&.each do |m|
      inherit_yr_level(m, y)
    end
  when Hash
    inherit_yr_level(json['asn:hasLevel'], y)
  end
end

def match_predicates_json_ld(json)
  json_keys = @asnjson.dig(id(json), :keys)
  jsonld_keys = json.keys.select { |x| x.include?(':') }
                    .reject { |x| %w[gem:hasChild gem:isChildOf dc:isPartOf esa:nominalYearLevel].include?(x) }
  if jsonld_keys.include?('dc:description')
    jsonld_keys = jsonld_keys.reject { |x| x == 'dc:description' }
    jsonld_keys.include?('dc:title') or jsonld_keys << 'dc:title'
  end
  jsonld_keys.sort!
  if json_keys != jsonld_keys && !jsonld_keys.include?('skos:prefLabel') # prefLabel means this is not a node, but an alias
    l = json['asn:contentdesc']
    l.is_a?(Hash) and l = json.dig('asn:statementLabel', 'literal')
    keys_diff = keys_diff(json_keys, jsonld_keys)
    @err.puts "#{@f}: JSON (#{l}): inconsistent predicates between JSON and JSON-LD (#{keys_diff}) in #{id(json)}"
  end
  json['gem:hasChild']&.each do |m|
    match_predicates_json_ld(m)
  end
  case json['asn:hasLevel']
  when Array
    json['asn:hasLevel']&.each do |m|
      match_predicates_json_ld(m)
    end
  when Hash
    match_predicates_json_ld(json['asn:hasLevel'])
  end
end

def keys_diff(arr1, arr2)
  ret = []
  arr1.each do |k|
    arr2.include?(k) or ret << "-#{k}"
  end
  arr2.each do |k|
    arr1.include?(k) or ret << "+#{k}"
  end
  ret.join(', ')
end

def extract_yr_levels(json)
  json.nil? and return []
  case json
  when Array
    json.map { |j| extract_yr_levels(j) }.flatten
  when Hash
    [id(json).sub(%r{^.+/}, '').sub("Foundation Year", "0").sub(/^Year\s+/, "")]
  else [json.sub(%r{^.+/}, '').sub("Foundation Year", "0").sub(/^Year\s+/, "")]
  end
end

def subset?(arr_a, arr_b)
  arr_a.sort!
  arr_b.sort!
  arr_a.select.with_index do |a, index|
    arr_b[index] == a
  end == arr_b
end

def validate_file(json)
  child_parent(json)
  namespaced_keys(json, nil)
  link_types(json)
  scot(json)
  scot_mappings(json)
  match_predicates_json_ld(json)
  year_level(json)
end

def read_file(f)
  JSON.parse(File.read(f))
end

def read_scot
  scot = JSON.parse(File.read(SCOT_FILE))
  @scotids = scot_id_parse(scot, {})
end

def read_scot_mapping
  f = CSV.parse(File.read(SCOT_MAPPING_FILE), headers: true)
  @mappings = {}
  f.each do |k|
    @mappings[k["mrac"]] ||= []
    @mappings[k["mrac"]] << "http://vocabulary.curriculum.edu.au/scot/#{k["scot"]}"
  end
end

def validate
  warn 'Processing...'
  @err = File.open('err.txt', 'w')
  read_scot
  read_scot_mapping
  @ids = {}
  @asnjson = {}
  Dir.glob("#{JSON_LD_FOLDER}/*.json") do |f|
    @f = f
    a = read_file(f)
    @ids = asn_jsonld_id_parse(a, @ids)
  end
  Dir.glob("#{JSON_FOLDER}/*.json") do |f|
    @f = f
    f == "#{JSON_FOLDER}/asn-node.json" and next
    a = read_file(f)
    @asnjson = asn_json_id_parse(a, @asnjson)
  end
  Dir.glob("#{JSON_LD_FOLDER}/*.json") do |f|
    @f = f
    a = read_file(f)
    validate_file(a)
  end
  warn '...Done'
  @err.close
end

validate
