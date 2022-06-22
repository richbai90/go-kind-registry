#! /bin/bash
set -o errexit

if [ -z ${reg_name} ]; then
    reg_name='kr'
fi
if [ -z ${reg_version} ]; then
    reg_version=1.24.0
fi

# Make sure the staging and outfile variables are set correctly
[ -z ${staging_dir} ] && echo "Staging Directory Not Set" >&2 && exit 1
[ -z ${outfile} ] && echo "Outifile Not Specified" >&2 && exit 1

# Make sure the staging directory exists
echo "Creating staging directory"
mkdir -p "$staging_dir"
# Save the registry image
echo "Saving the registry image"
docker save -o ${staging_dir}/registry.tar registry:2
# Save the kind image
echo "Saving the kind image"
docker save -o ${staging_dir}/kind.tar kindest/node:v${reg_version}
# Pull alpine
echo "Pulling alpine"
docker pull alpine:3
# Save the alpine image
echo "Saving the alpine image"
docker save -o ${staging_dir}/alpine.tar alpine:3
# Export the volumes
echo "Exporting volumes"
docker run --rm --volumes-from ${reg_name} -v ${staging_dir}:/backup alpine:3 \
    sh -c "tar -cvzf /backup/volumes.tar.gz /var/lib/registry /root"
# Complete the bundle
echo "Compiling bundle"
docker run --rm -v ${staging_dir}:/backup alpine:3 \
    sh -c "tar -cvzf /backup/bundle-complete.tar.gz /backup/registry.tar /backup/kind.tar /backup/alpine.tar /backup/volumes.tar.gz"
# Permissions get funky after mounting in container running as root
echo "Correcting permissions"
chmod -R 0760 "$staging_dir"

# move the output bundle
echo "Exporting Bundle"
mv ${staging_dir}/bundle-complete.tar.gz ${outfile}

#cleanup
echo "Cleaning Up"
rm -rf ${staging_dir}
