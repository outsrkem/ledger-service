<template>
  <el-pagination
    @size-change="handleSizeChange"
    @current-change="handleCurrentChange"
    :current-page="currentPage"
    :page-size="pageSize"
    :page-sizes="pageSizes"
    layout="total, sizes, prev, pager, next, jumper"
    :total="pageTotal"
  />
</template>

<script>
export default {
    name: "PaginationIndex",
    emits: {
    // Synchronize pagination parameters after component mounted: (size, page)
        syncsize: (size, page) => {
            return typeof size === "number" && typeof page === "number" && page >= 1;
        },
        // Triggered when page size changes: (page_size)
        SizeChange: (page_size) => {
            return typeof page_size === "number" && page_size > 0;
        },
        // Triggered when current page number changes: (current_page, pageSize)
        CurrentChange: (current_page, pageSize) => {
            return (
                typeof current_page === "number" &&
        current_page >= 1 &&
        typeof pageSize === "number"
            );
        },
    },
    props: {
        pageTotal: { type: Number, default: 0 }, // Total number of items
        pageSize: { type: Number, default: 15 }, // Current page size
        autoSync: { type: Boolean, default: true }, // Auto sync pagination params from URL on mounted
    },
    data() {
        return {
            // Current active page
            currentPage: 1,

            // Available page size options
            pageSizes: [10, 15, 25, 50, 100],

            // Avoid repeated triggering of sync event
            isSyncing: false,
        };
    },
    mounted() {
        if (!this.autoSync) return;

        // Read pagination parameters from URL query on component mount
        if (this.isSyncing) return;
        this.isSyncing = true;
        try {
            const { p, s } = this.$route.query;
            let validSize = this.pageSize;
            if (s) {
                const numS = Number(s);
                if (this.pageSizes.includes(numS)) {
                    validSize = numS;
                }
            }
            if (p) {
                const po = Number(p);
                if (!Number.isNaN(po) && po >= 1) {
                    this.currentPage = po;
                }
            }
            this.$emit("syncsize", validSize, this.currentPage);
        } finally {
            this.isSyncing = false;
        }
    },
    methods: {
    /**
     * Handle page size change event
     * @param {number} page_size - The new page size
     */
        handleSizeChange(page_size) {
            this.currentPage = 1; // Reset to first page when size changes
            this.$emit("SizeChange", page_size);
            this.setUrlQuery(1, page_size);
        },

        /**
     * Handle current page change event
     * @param {number} current_page - The new current page number
     */
        handleCurrentChange(current_page) {
            this.currentPage = current_page;
            this.$emit("CurrentChange", current_page, this.pageSize);
            this.setUrlQuery(current_page, this.pageSize);
        },

        /**
     * Update URL query parameters for pagination (s first, then p)
     * @param {number} p - Page number
     * @param {number} s - Page size
     */
        setUrlQuery(p = this.currentPage, s = this.pageSize) {
            let finalS = this.pageSizes.includes(Number(s)) ? s : this.pageSizes[1];
            this.$router.push({
                path: this.$route.path,
                query: {
                    ...this.$route.query,
                    s: finalS,
                    p,
                },
            });
        },

        /**
     * Navigate with pagination parameters, supports both 'name' and 'path' route formats
     * @param {Object} routeObj - Route object {name:'xxx', params:{}} or {path:'xxx'}
     * @param {Object} extraQuery - Additional query parameters (pagination s/p will be merged automatically)
     */
        goWithPageParams(routeObj, extraQuery = {}) {
            const { p, s } = this.$route.query;
            const numS = Number(s);
            const safeS = this.pageSizes.includes(numS) ? numS : null;
            const pageQuery = {
                ...(safeS ? { s: safeS } : {}),
                ...(p ? { p } : {}),
                ...extraQuery,
            };
            const targetRoute = {
                ...routeObj,
                query: pageQuery,
            };
            this.$router.push(targetRoute);
        },

        /**
     * Navigate back to a list page while preserving pagination parameters
     * @param {string} listPath - The path of the list page to navigate to
     */
        backToPageList(listPath) {
            const { p, s } = this.$route.query;
            const numS = Number(s);
            const safeS = this.pageSizes.includes(numS) ? numS : null;
            this.$router.push({
                path: listPath,
                query: {
                    ...(safeS && { s: safeS }),
                    ...(p && { p }),
                },
            });
        },
    },
};
</script>
