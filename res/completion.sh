__veidemannctl_parse_get() {
    declare -rA tables=( ["entity"]="config_crawl_entities" ["browser"]="config_browser_configs" \
      ["crawlconfig"]="config_crawl_configs" ["group"]="config_crawl_host_group_configs" ["job"]="config_crawl_jobs" \
      ["politeness"]="config_politeness_configs" ["role"]="dog" ["schedule"]="config_crawl_schedule_configs" \
      ["script"]="config_browser_scripts" ["seed"]="config_seeds" )
    local table=${tables[$1]} query template veidemannctl_out

	if [ -z "$table" ]; then
      return 0
    fi

    query="r.table('${table}')"
	if [ -n "$cur" ]; then
      query="${query}.between('${cur}', '${cur}z', {index:'id'})"
    fi
	query="${query}.orderBy({index:'id'}).pluck('id')"
    template="{{println .id}}"
    if mapfile -t veidemannctl_out < <( veidemannctl report query "${query}" -o template -t"${template}" -s20 2>/dev/null ); then
        mapfile -t COMPREPLY < <( compgen -W "$( printf '%q ' "${veidemannctl_out[@]}" )" -- "$cur" | awk '/ / { print "\""$0"\"" } /^[^ ]+$/ { print $0 }' )
    fi
}

__veidemannctl_get_resource() {
    if [[ ${#nouns[@]} -eq 0 ]]; then
        return 1
    fi
    __veidemannctl_parse_get ${nouns[${#nouns[@]} -1]}
    if [[ $? -eq 0 ]]; then
        return 0
    fi
}

__veidemannctl_query_resource() {
    local veidemannctl_out
    if mapfile -t veidemannctl_out < <( veidemannctl report query -q 2>/dev/null ); then
        mapfile -t COMPREPLY < <( compgen -W "$( printf '%q ' "${veidemannctl_out[@]}" )" -- "$cur" | awk '/ / { print "\""$0"\"" } /^[^ ]+$/ { print $0 }' )
    fi
}

__veidemannctl_custom_func() {
    case ${last_command} in
        veidemannctl_get)
            __veidemannctl_get_resource
            return
            ;;
        veidemannctl_report_query)
			__veidemannctl_query_resource
			return
			;;
        *)
            ;;
    esac
}

__veidemannctl_get_name() {
    if [[ ${#nouns[@]} -eq 0 ]]; then
        return 1
    fi
	local noun
	noun=${nouns[${#nouns[@]} -1]}
    local template
    template="{{println .meta.name}}"
    local veidemannctl_out
    if mapfile -t veidemannctl_out < <( veidemannctl get "$noun" -n "^${cur}" -s20 -o template -t "${template}" 2>/dev/null ); then
	    mapfile -t COMPREPLY < <( printf '%q\n' "${veidemannctl_out[@]}" | awk -v IGNORECASE=1 -v p="$cur" '{p==substr($0,0,length(p))} { print $0 }' )
    fi
}
